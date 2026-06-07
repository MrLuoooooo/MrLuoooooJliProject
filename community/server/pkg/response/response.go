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
	CodeSuccess      = 0
	CodeInvalidParam = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerBusy   = 500
	// 401xxx 鉴权
	CodeTokenExpired = 401001
	CodeTokenInvalid = 401002
	// 400xxx 参数/业务校验
	CodeUserExists       = 400001
	CodeEmailExists      = 400002
	CodeWrongPassword    = 400003
	CodeSelfMessage      = 400004
	CodeUserNotFound     = 400005
	CodeAlreadyLiked     = 400006
	CodeAlreadyFavorited = 400007
	CodeAlreadyFollowing = 400008
	CodeNotFollowing     = 400009
	CodeSelfFollow       = 400010
	CodeCommentNotFound  = 400011
	// 403xxx 权限
	CodeNotOwner       = 403001
	CodeAdminRequired  = 403002
	CodeUserBanned     = 403003
	// 429 限流
	CodeTooManyRequests = 429
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
		CodeSuccess:      "ok",
		CodeInvalidParam: "参数错误",
		CodeUnauthorized: "未登录",
		CodeForbidden:    "无权限",
		CodeNotFound:     "资源不存在",
		CodeServerBusy:   "服务器繁忙",
		// 401xxx
		CodeTokenExpired: "令牌已过期",
		CodeTokenInvalid: "令牌无效",
		// 400xxx
		CodeUserExists:       "用户名已存在",
		CodeEmailExists:      "邮箱已注册",
		CodeWrongPassword:    "密码错误",
		CodeSelfMessage:      "不能给自己发消息",
		CodeUserNotFound:     "用户不存在",
		CodeAlreadyLiked:     "已点赞",
		CodeAlreadyFavorited: "已收藏",
		CodeAlreadyFollowing: "已关注",
		CodeNotFollowing:     "未关注",
		CodeSelfFollow:       "不能关注自己",
		CodeCommentNotFound:  "评论不存在",
		// 403xxx
		CodeNotOwner:      "非本人操作",
		CodeAdminRequired: "需要管理员权限",
		CodeUserBanned:    "账号已被禁用",
	}
	if msg, ok := msgs[code]; ok {
		return msg
	}
	return "未知错误"
}
