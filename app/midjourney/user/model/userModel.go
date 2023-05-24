package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	_                              UserModel = (*customUserModel)(nil)
	userRegisterRowsExpectAutoSet            = strings.Join(stringx.Remove(userFieldNames, "`id`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	cacheMidjourneyUserEmailPrefix           = "cache:midjourney:user:email:"
)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindOneByEmail(ctx context.Context, email string) (*User, error)
		RegisterByEmail(ctx context.Context, data *User) (sql.Result, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c),
	}
}

func (m *defaultUserModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	midjourneyUserEmailKey := fmt.Sprintf("%s%v", cacheMidjourneyUserEmailPrefix, email)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, midjourneyUserEmailKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, email); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) RegisterByEmail(ctx context.Context, data *User) (sql.Result, error) {
	midjourneyUserIdKey := fmt.Sprintf("%s%v", cacheMidjourneyUserIdPrefix, data.Id)
	midjourneyUserEmailKey := fmt.Sprintf("%s%v", cacheMidjourneyUserEmailPrefix, data.Email)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, userRegisterRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.Password, data.Avatar, data.Phone, data.Email, data.State, data.CreatedAt)
	}, midjourneyUserEmailKey, midjourneyUserIdKey)
	return ret, err
}
