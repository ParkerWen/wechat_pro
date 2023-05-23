package logic

import (
	"fmt"

	"github.com/ParkerWen/wechat_pro/app/mqueue/job/jobtype"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

func (s *MqueueScheduler) updateMJRecordScheduler() {
	task := asynq.NewTask(jobtype.ScheduleUpdateMJRecord, nil)

	entryID, err := s.svcCtx.Scheduler.Register("*/1 * * * *", task)
	if err != nil {
		logx.WithContext(s.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【settleRecordScheduler】 registered  err:%+v , task:%+v", err, task)
	}
	fmt.Printf("【settleRecordScheduler】 registered an  entry: %q \n", entryID)
}
