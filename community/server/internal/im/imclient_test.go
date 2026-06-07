package im

import (
	"testing"
)

// TestIMClientInterface 编译期检查：*client 实现了 IMClient 接口
// 附带验证 UserIDToStr 工具函数的输入输出
func TestIMClientInterface(t *testing.T) {
	// 编译期检查：*client 实现了 IMClient
	var _ IMClient = (*client)(nil)

	// 检查 UserIDToStr 工具函数
	tests := []struct {
		input uint
		want  string
	}{
		{1, "1"},
		{100, "100"},
		{999999, "999999"},
	}

	for _, tt := range tests {
		got := UserIDToStr(tt.input)
		if got != tt.want {
			t.Errorf("UserIDToStr(%d) = %s; 期望 %s", tt.input, got, tt.want)
		}
	}
}

// TestGenerateSignature 验证签名生成算法
func TestGenerateSignature(t *testing.T) {
	c := &client{
		appSecret: "test-secret",
		appKey:    "test-key",
		baseURL:   "http://localhost:9001",
	}

	// 已知输入输出
	sig := c.generateSignature("123456", "1700000000000")
	if sig == "" {
		t.Error("签名不应为空")
	}
	if len(sig) != 40 {
		t.Errorf("SHA1 签名长度应为 40，实际为 %d", len(sig))
	}
}

// TestBuildSignatureHeaders 验证签名头里五个必需字段都存在
func TestBuildSignatureHeaders(t *testing.T) {
	c := &client{
		appSecret: "secret",
		appKey:    "key",
		baseURL:   "http://localhost:9001",
	}

	headers := c.buildSignatureHeaders()
	requiredKeys := []string{"Content-Type", "appkey", "nonce", "timestamp", "signature"}
	for _, key := range requiredKeys {
		if _, ok := headers[key]; !ok {
			t.Errorf("缺少请求头: %s", key)
		}
	}

	if headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type 应为 application/json")
	}
	if headers["appkey"] != "key" {
		t.Errorf("appkey 应为 key")
	}
}
