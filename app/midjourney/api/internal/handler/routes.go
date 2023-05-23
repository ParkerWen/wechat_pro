// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/fetchTask",
				Handler: fetchTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/listTask",
				Handler: listTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/imagine",
				Handler: imagineHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/upscale",
				Handler: upscaleHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/variation",
				Handler: variationHandler(serverCtx),
			},
		},
		rest.WithPrefix("/mj/v1"),
	)
}
