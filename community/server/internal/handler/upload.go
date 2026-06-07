package handler

import (
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService UploadService
}

func NewUploadHandler(uploadService UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Tags 文件
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} response.Response
// @Router /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "请选择文件")
		return
	}

	if file.Size > 10<<20 {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "文件大小不能超过 10MB")
		return
	}

	result, err := h.uploadService.Upload(file)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}
