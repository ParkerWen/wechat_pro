package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/ParkerWen/wechat_pro/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	cacheKeyPrefix = "gowechat_work_"
	uploadTempFile = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
)

type UploadTempFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadTempFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTempFileLogic {
	return &UploadTempFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadTempFileLogic) UploadTempFile(req *types.UploadTempFileReq) (*types.UploadTempFileResp, error) {
	res, err := http.Get(req.Media)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWrite, err := bodyWriter.CreateFormFile("media", tool.Krand(10, tool.KC_RAND_KIND_ALL))
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	len, err := fileWrite.Write(fileBytes)
	if err != nil {
		return nil, err
	}
	if len <= 0 {
		return nil, err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// 上传临时素材
	resp, err := http.Post(fmt.Sprintf(uploadTempFile, req.AccessToken, "image"), contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析返回的数据流
	var uploadResp types.UploadTempFileResp
	err = json.Unmarshal(respBody, &uploadResp)
	if err != nil {
		return nil, err
	}
	return &uploadResp, nil
}
