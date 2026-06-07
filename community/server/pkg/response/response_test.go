package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func performRequest(fn func(c *gin.Context)) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fn(c)
	return w
}

func TestSuccess(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Success(c, gin.H{"key": "value"})
	})

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != CodeSuccess {
		t.Errorf("期望 code %d，实际 %d", CodeSuccess, resp.Code)
	}
}

func TestErrorWithMsg(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		ErrorWithMsg(c, CodeInvalidParam, "参数错误")
	})

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != CodeInvalidParam {
		t.Errorf("期望 code %d，实际 %d", CodeInvalidParam, resp.Code)
	}
	if resp.Msg != "参数错误" {
		t.Errorf("期望 msg '参数错误'，实际 %s", resp.Msg)
	}
}

func TestUnauthorized(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Error(c, CodeUnauthorized)
	})

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != CodeUnauthorized {
		t.Errorf("期望 code %d，实际 %d", CodeUnauthorized, resp.Code)
	}
}

func TestNotFound(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Error(c, CodeNotFound)
	})

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != CodeNotFound {
		t.Errorf("期望 code %d，实际 %d", CodeNotFound, resp.Code)
	}
}
