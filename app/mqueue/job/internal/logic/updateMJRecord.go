package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/ParkerWen/wechat_pro/app/mqueue/job/internal/svc"
	"github.com/hibiken/asynq"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMJRecordHandler struct {
	svcCtx *svc.ServiceContext
}

type TaskFetchRes struct {
	TaskId      string `db:"id"`
	Action      string `db:"action"`
	Prompt      string `db:"prompt"`
	PromptEn    string `db:"promptEn"`
	ImageUrl    string `db:"imageUrl"`
	Description string `db:"description"`
	Status      string `db:"status"`
	State       string `db:"state"`
	SubmitTime  int64  `db:"submitTime"`
	FinishTime  int64  `db:"finishTime"`
}

func NewUpdateMJRecordHandler(svcCtx *svc.ServiceContext) *UpdateMJRecordHandler {
	return &UpdateMJRecordHandler{
		svcCtx: svcCtx,
	}
}

func (h *UpdateMJRecordHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	fmt.Printf("shcedule job UpdateByMJ -----> every one minute exec \n")

	whereBuilder := h.svcCtx.TaskModel.RowBuilder()
	list, err := h.svcCtx.TaskModel.FindByStatus(ctx, "PENDING", whereBuilder)
	if err != nil && err != model.ErrNotFound {
		return errors.Wrapf(errors.New("数据库繁忙,请稍后再试"), "Failed to get pending task err : %v ", err)
	}
	if len(list) > 0 {
		for _, taskInfo := range list {
			request := gorequest.New()
			resp, body, _ := request.Get(fmt.Sprintf("http://38.95.233.164:8080/mj/task/%s/fetch", taskInfo.TaskId)).End()
			defer resp.Body.Close()
			if len(body) <= 0 {
				continue
			}
			var taskFetchRes TaskFetchRes
			err := json.NewDecoder(resp.Body).Decode(&taskFetchRes)
			if err != nil {
				logx.Error(err)
				continue
			}
			if taskFetchRes.Status == "SUCCESS" {
				// UpdateByMJ
				taskInfo.ImageUrl = taskFetchRes.ImageUrl
				taskInfo.Description = taskFetchRes.Description
				taskInfo.Status = taskFetchRes.Status
				taskInfo.UpdatedAt = time.Now().Unix()
				err = h.svcCtx.TaskModel.UpdateByMJ(ctx, taskInfo)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
