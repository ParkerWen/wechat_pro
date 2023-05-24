package logic

import (
	"context"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/types"
	"github.com/ParkerWen/wechat_pro/common/ctxdata"
	"github.com/ParkerWen/wechat_pro/common/xerr"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(req *types.GenerateTokenReq) (*types.GenerateTokenResp, error) {
	now := time.Now().Unix()
	AccessSecret := l.svcCtx.Config.Auth.AccessSecret
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(AccessSecret, now, accessExpire, req.UserId)
	if err != nil {
		logx.Error(err)
		return nil, errors.Wrapf(xerr.NewErrMsg("生成token失败"), "GetJwtToken Err UserId: %d , Err: %v", req.UserId, err)
	}

	return &types.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[ctxdata.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
