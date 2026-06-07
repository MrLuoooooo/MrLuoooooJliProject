package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}
	if token == "" {
		t.Fatal("token 不应为空")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken 失败: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("期望 UserID 1，实际 %d", claims.UserID)
	}
	if claims.UserName != "testuser" {
		t.Errorf("期望 UserName testuser，实际 %s", claims.UserName)
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.here")
	if err == nil {
		t.Error("无效 token 应返回错误")
	}
}

func TestParseExpiredToken(t *testing.T) {
	claims := Claims{
		UserID:   1,
		UserName: "test",
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "community-server",
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := ParseToken(tokenStr)
	if err == nil {
		t.Error("过期 token 应返回错误")
	}
}

func TestRefreshToken(t *testing.T) {
	token, err := GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	newToken, err := RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken 失败: %v", err)
	}

	// 验证新 token 包含正确的 claims
	claims, err := ParseToken(newToken)
	if err != nil {
		t.Fatalf("新 token 解析失败: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("期望 UserID 1，实际 %d", claims.UserID)
	}
	if claims.UserName != "testuser" {
		t.Errorf("期望 UserName testuser，实际 %s", claims.UserName)
	}
}
