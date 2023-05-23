package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/mqueue/job/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/mqueue/job/jobtype"
	"github.com/hibiken/asynq"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register Job
func (c *CronJob) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	// Scheduler Job
	mux.Handle(jobtype.ScheduleUpdateMJRecord, NewUpdateMJRecordHandler(c.svcCtx))

	return mux
}
