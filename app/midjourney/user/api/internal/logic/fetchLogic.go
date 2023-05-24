package logic

import (
	"context"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/types"
	"github.com/ParkerWen/wechat_pro/common/ctxdata"
	"github.com/ParkerWen/wechat_pro/common/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchLogic {
	return &FetchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchLogic) Fetch(req *types.FetchReq) (*types.FetchResp, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Fetch User DB ERROR , id: %d , err: %v", userId, err)
	}
	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "id: %d", userId)
	}
	data := make(map[string]any)
	data["id"] = user.Id
	data["name"] = user.Name
	data["avatar"] = user.Avatar
	data["email"] = user.Email
	data["phone"] = user.Phone
	return &types.FetchResp{
		Code: http.StatusOK,
		Msg:  "SUCCESS",
		Data: data,
	}, nil
}
