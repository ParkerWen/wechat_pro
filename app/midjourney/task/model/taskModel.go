package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/ParkerWen/wechat_pro/common/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ TaskModel = (*customTaskModel)(nil)

var (
	taskImagineRowsExpectAutoSet   = strings.Join(stringx.Remove(taskFieldNames, "`id`"), ",")
	taskImagineRowsWithPlaceHolder = strings.Join(stringx.Remove(taskFieldNames, "`id`", "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`"), "=?,") + "=?"
)

type (
	// TaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskModel.
	TaskModel interface {
		taskModel
		RowBuilder() squirrel.SelectBuilder
		InsertByImagine(ctx context.Context, data *Task) (sql.Result, error)
		UpdateByMJ(ctx context.Context, newData *Task) error
		FindByStatus(ctx context.Context, status string, rowBuilder squirrel.SelectBuilder) ([]*Task, error)
		FindByParentTaskIdAndAction(ctx context.Context, parent_task_id, action string, rowBuilder squirrel.SelectBuilder) ([]*Task, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Task, error)
		FindPageListByIdDesc(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Task, error)
	}

	customTaskModel struct {
		*defaultTaskModel
	}
)

// NewTaskModel returns a model for the database table.
func NewTaskModel(conn sqlx.SqlConn, c cache.CacheConf) TaskModel {
	return &customTaskModel{
		defaultTaskModel: newTaskModel(conn, c),
	}
}

func (m *defaultTaskModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(taskRows).From(m.table)
}

func (m *defaultTaskModel) InsertByImagine(ctx context.Context, data *Task) (sql.Result, error) {
	midjourneyTaskIdKey := fmt.Sprintf("%s%v", cacheMidjourneyTaskIdPrefix, data.Id)
	midjourneyTaskTaskIdKey := fmt.Sprintf("%s%v", cacheMidjourneyTaskTaskIdPrefix, data.TaskId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, taskImagineRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.TaskId, data.UserId, data.ParentTaskId, data.Action, data.Index, data.Prompt, data.ImageUrl, data.Description, data.Status, data.State, data.CreatedAt, data.UpdatedAt)
	}, midjourneyTaskIdKey, midjourneyTaskTaskIdKey)
	return ret, err
}

func (m *defaultTaskModel) UpdateByMJ(ctx context.Context, newData *Task) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	midjourneyTaskIdKey := fmt.Sprintf("%s%v", cacheMidjourneyTaskIdPrefix, data.Id)
	midjourneyTaskTaskIdKey := fmt.Sprintf("%s%v", cacheMidjourneyTaskTaskIdPrefix, data.TaskId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, taskImagineRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.TaskId, newData.ParentTaskId, newData.Action, newData.Index, newData.Prompt, newData.ImageUrl, newData.Description, newData.Status, newData.State, newData.UpdatedAt, newData.Id)
	}, midjourneyTaskIdKey, midjourneyTaskTaskIdKey)
	return err
}

func (m *defaultTaskModel) FindByStatus(ctx context.Context, status string, rowBuilder squirrel.SelectBuilder) ([]*Task, error) {
	query, values, err := rowBuilder.Where("status = ?", status).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []*Task
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultTaskModel) FindByParentTaskIdAndAction(ctx context.Context, parent_task_id, action string, rowBuilder squirrel.SelectBuilder) ([]*Task, error) {
	query, values, err := rowBuilder.Where("parent_task_id = ?", parent_task_id).Where("action = ?", action).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []*Task
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultTaskModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Task, error) {
	userId := ctxdata.GetUidFromCtx(ctx)
	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("state = ?", "valid").Where("status = ?", "SUCCESS").Where("user_id = ?", userId).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Task
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultTaskModel) FindPageListByIdDesc(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Task, error) {
	userId := ctxdata.GetUidFromCtx(ctx)
	if preMinId > 0 {
		rowBuilder = rowBuilder.Where("id < ?", preMinId)
	}
	query, values, err := rowBuilder.Where("state = ?", "valid").Where("status = ?", "SUCCESS").Where("user_id = ?", userId).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []*Task
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
