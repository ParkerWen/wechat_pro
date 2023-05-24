package logic

import (
	"context"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/model"
	"github.com/ParkerWen/wechat_pro/common/tool"
	"github.com/ParkerWen/wechat_pro/common/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByEmailLogic {
	return &LoginByEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByEmailLogic) LoginByEmail(req *types.LoginByEmailReq) (*types.LoginByEmailResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "根据邮箱查询用户信息失败，Email: %s, Err: %v", req.Email, err)
	}
	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("该邮箱还未注册"), "Email: %s", req.Email)
	}
	if !(tool.Md5ByString(req.Password) == user.Password) {
		return nil, errors.Wrap(xerr.NewErrMsg("账号或密码不正确"), "密码匹配出错")
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	generateTokenResp, err := generateTokenLogic.GenerateToken(&types.GenerateTokenReq{
		UserId: user.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("生成token失败了"), "GenerateToken UserId : %d", user.Id)
	}
	var data = make(map[string]interface{})
	data["access_token"] = generateTokenResp.AccessToken
	data["access_expire"] = generateTokenResp.AccessExpire
	data["refresh_after"] = generateTokenResp.RefreshAfter
	return &types.LoginByEmailResp{
		Code: http.StatusOK,
		Msg:  "SUCCESS",
		Data: data,
	}, nil
}
