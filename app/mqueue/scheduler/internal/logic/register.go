package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/mqueue/scheduler/internal/svc"
)

type MqueueScheduler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCornScheduler(ctx context.Context, svcCtx *svc.ServiceContext) *MqueueScheduler {
	return &MqueueScheduler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (s *MqueueScheduler) Register() {
	s.updateMJRecordScheduler()
}
