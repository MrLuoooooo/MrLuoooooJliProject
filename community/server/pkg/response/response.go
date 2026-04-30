package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess         = 0
	CodeInvalidParam    = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeTooManyRequests = 429
	CodeServerBusy      = 500
	CodeTokenExpired    = 401001
	CodeTokenInvalid    = 401002
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  getMsgByCode(code),
	})
}

func ErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

func ErrorWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  getMsgByCode(code),
		Data: data,
	})
}

func getMsgByCode(code int) string {
	msgs := map[int]string{
		CodeSuccess:      "success",
		CodeInvalidParam: "参数错误",
		CodeUnauthorized: "未授权",
		CodeForbidden:    "禁止访问",
		CodeNotFound:     "资源不存在",
		CodeServerBusy:   "服务器繁忙",
		CodeTokenExpired: "令牌已过期",
		CodeTokenInvalid: "令牌无效",
	}
	if msg, ok := msgs[code]; ok {
		return msg
	}
	return "unknown error"
}
