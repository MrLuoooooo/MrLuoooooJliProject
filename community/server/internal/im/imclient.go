package im

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	appKey     string
	appSecret  string
	httpClient *http.Client
}

func NewClient(baseURL, appKey, appSecret string) *Client {
	return &Client{
		baseURL:    baseURL,
		appKey:     appKey,
		appSecret:  appSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) buildSignatureHeaders() map[string]string {
	nonce := fmt.Sprintf("%d", rand.Int31n(1000000))
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := c.generateSignature(nonce, timestamp)

	return map[string]string{
		"Content-Type": "application/json",
		"appkey":       c.appKey,
		"nonce":        nonce,
		"timestamp":    timestamp,
		"signature":    signature,
	}
}

func (c *Client) generateSignature(nonce, timestamp string) string {
	str := fmt.Sprintf("%s%s%s", c.appSecret, nonce, timestamp)
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Client) Do(method, path string, body interface{}) ([]byte, error) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	headers := c.buildSignatureHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) RegisterUser(userId, nickname string) ([]byte, error) {
	body := map[string]interface{}{
		"user_id":  userId,
		"nickname": nickname,
	}
	return c.Do("POST", "/apigateway/users/register", body)
}

func (c *Client) SendPrivateMsg(senderId, targetId, msgType, content string) ([]byte, error) {
	body := map[string]interface{}{
		"sender_id":   senderId,
		"target_id":   targetId,
		"msg_type":    msgType,
		"msg_content": content,
	}
	return c.Do("POST", "/apigateway/messages/private/send", body)
}

func (c *Client) QueryUserInfo(userId string) ([]byte, error) {
	return c.Do("GET", fmt.Sprintf("/apigateway/users/info?user_id=%s", userId), nil)
}

func (c *Client) CreateGroup(groupId, groupName, ownerId string, memberIds []string) ([]byte, error) {
	body := map[string]interface{}{
		"group_id":   groupId,
		"group_name": groupName,
		"owner_id":   ownerId,
		"member_ids": memberIds,
	}
	return c.Do("POST", "/apigateway/groups/add", body)
}
