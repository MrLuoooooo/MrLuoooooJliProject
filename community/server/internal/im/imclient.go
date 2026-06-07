package im

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"community-server/internal/config"
)

// IMClient 定义了与 JuggleIM 服务交互的接口
type IMClient interface {
	// RegisterUser 在 IM 系统中注册用户
	RegisterUser(userID string, nickname string) error
	// SendPrivateMsg 发送私聊消息（实时推送）
	SendPrivateMsg(senderID, targetID, content string) error
	// SendSystemMsg 发送系统消息
	SendSystemMsg(senderID, targetID, content string) error
	// SendGroupMsg 发送群聊消息
	SendGroupMsg(senderID, groupID, content string) error
	// QueryUserInfo 查询 IM 用户信息
	QueryUserInfo(userID string) (map[string]interface{}, error)
	// CreateGroup 创建群组
	CreateGroup(groupID, groupName, ownerID string, memberIDs []string) error
	// QueryOnlineStatus 查询用户在线状态
	QueryOnlineStatus(userIDs []string) (map[string]bool, error)
	// SendBroadcastMsg 发送全站广播消息
	SendBroadcastMsg(senderID, content string) error
	// AddBot 注册机器人
	AddBot(botID, nickname, webhookURL string) error
}

type client struct {
	baseURL    string
	appKey     string
	appSecret  string
	httpClient *http.Client
}

func NewIMClient(cfg *config.Config) IMClient {
	imCfg := cfg.IM
	return &client{
		baseURL:   imCfg.BaseURL,
		appKey:    imCfg.AppKey,
		appSecret: imCfg.AppSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) buildSignatureHeaders() map[string]string {
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

func (c *client) generateSignature(nonce, timestamp string) string {
	str := fmt.Sprintf("%s%s%s", c.appSecret, nonce, timestamp)
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *client) do(method, path string, body interface{}) ([]byte, error) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	headers := c.buildSignatureHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// apiResponse JuggleIM API 通用响应结构
type apiResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func (c *client) RegisterUser(userID, nickname string) error {
	body := map[string]interface{}{
		"user_id":  userID,
		"nickname": nickname,
	}
	resp, err := c.do("POST", "/apigateway/users/register", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM注册失败: %s", r.Msg)
	}
	return nil
}

func (c *client) SendPrivateMsg(senderID, targetID, content string) error {
	body := map[string]interface{}{
		"sender_id":   senderID,
		"target_id":   targetID,
		"msg_type":    "jg:text",
		"msg_content": content,
	}
	resp, err := c.do("POST", "/apigateway/messages/private/send", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM发送消息失败: %s", r.Msg)
	}
	return nil
}

func (c *client) SendSystemMsg(senderID, targetID, content string) error {
	body := map[string]interface{}{
		"sender_id":   senderID,
		"target_id":   targetID,
		"msg_type":    "jg:text",
		"msg_content": content,
	}
	resp, err := c.do("POST", "/apigateway/messages/system/send", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM发送系统消息失败: %s", r.Msg)
	}
	return nil
}

func (c *client) SendGroupMsg(senderID, groupID, content string) error {
	body := map[string]interface{}{
		"sender_id":   senderID,
		"target_id":   groupID,
		"msg_type":    "jg:text",
		"msg_content": content,
	}
	resp, err := c.do("POST", "/apigateway/messages/group/send", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM发送群消息失败: %s", r.Msg)
	}
	return nil
}

func (c *client) QueryUserInfo(userID string) (map[string]interface{}, error) {
	resp, err := c.do("GET", fmt.Sprintf("/apigateway/users/info?user_id=%s", userID), nil)
	if err != nil {
		return nil, err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("IM查询用户失败: %s", r.Msg)
	}
	return r.Data, nil
}

func (c *client) CreateGroup(groupID, groupName, ownerID string, memberIDs []string) error {
	body := map[string]interface{}{
		"group_id":   groupID,
		"group_name": groupName,
		"owner_id":   ownerID,
		"member_ids": memberIDs,
	}
	resp, err := c.do("POST", "/apigateway/groups/add", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM创建群组失败: %s", r.Msg)
	}
	return nil
}

// UserIDToStr 将 uint user_id 转为 JuggleIM 使用的字符串格式
func UserIDToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

type onlineStatusResp struct {
	Items []struct {
		UserID   string `json:"user_id"`
		IsOnline bool   `json:"is_online"`
	} `json:"items"`
}

func (c *client) QueryOnlineStatus(userIDs []string) (map[string]bool, error) {
	body := map[string]interface{}{
		"user_ids": userIDs,
	}
	resp, err := c.do("POST", "/apigateway/users/onlinestatus/query", body)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data onlineStatusResp `json:"data"`
	}
	if err := json.Unmarshal(resp, &raw); err != nil {
		return nil, fmt.Errorf("解析在线状态失败: %w", err)
	}
	if raw.Code != 0 {
		return nil, fmt.Errorf("IM查询在线状态失败: %s", raw.Msg)
	}

	result := make(map[string]bool, len(raw.Data.Items))
	for _, item := range raw.Data.Items {
		result[item.UserID] = item.IsOnline
	}
	return result, nil
}

func (c *client) SendBroadcastMsg(senderID, content string) error {
	body := map[string]interface{}{
		"sender_id":   senderID,
		"msg_type":    "jg:text",
		"msg_content": content,
	}
	resp, err := c.do("POST", "/apigateway/messages/broadcast/send", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM广播失败: %s", r.Msg)
	}
	return nil
}

func (c *client) AddBot(botID, nickname, webhookURL string) error {
	body := map[string]interface{}{
		"bot_id":   botID,
		"nickname": nickname,
		"bot_type": 1,
		"webhook":  webhookURL,
	}
	resp, err := c.do("POST", "/apigateway/bots/add", body)
	if err != nil {
		return err
	}
	var r apiResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return fmt.Errorf("IM添加机器人失败: %s", r.Msg)
	}
	return nil
}
