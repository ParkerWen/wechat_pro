package logic

import (
	"context"
	"net/http"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/model"
	"github.com/ParkerWen/wechat_pro/common/tool"
	"github.com/ParkerWen/wechat_pro/common/xerr"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterByEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterByEmailLogic {
	return &RegisterByEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterByEmailLogic) RegisterByEmail(req *types.RegisterByEmailReq) (*types.RegisterByEmailResp, error) {
	// 判断用户的输入
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	// 判断该邮箱是否已经被注册
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Email: %s, Err: %v", req.Email, err)
	}
	// 如果存在就抛出错误
	if user != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Email has been registered"), "Register user exists Email: %s, Err: %v", req.Email, err)
	}
	// 不存在则注册 （事务模式）
	if err := l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		user := new(model.User)
		user.Email = req.Email
		user.State = "valid"
		if len(user.Name) == 0 {
			user.Name = tool.Krand(8, tool.KC_RAND_KIND_ALL)
		}
		user.Password = tool.Md5ByString(req.Password)
		user.CreatedAt = time.Now().Unix()
		res, err := l.svcCtx.UserModel.RegisterByEmail(l.ctx, user)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err: %v, user: %+v", err, user)
		}
		_, err = res.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user insertResult.LastInsertId err: %v, user: %+v", err, user)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &types.RegisterByEmailResp{
		Code: http.StatusOK,
		Msg:  "SUCCESS",
	}, nil
}
