package service

import (
	"fmt"
	"sync"
)

// mockIMClient 是 IMClient 的测试替身，记录所有调用
type mockIMClient struct {
	mu sync.Mutex

	registeredUsers map[string]string        // userID -> nickname
	privateMsgs     []privateMsgCall         // 发送的私信
	systemMsgs      []systemMsgCall          // 发送的系统消息
	groupMsgs       []string                 // 群聊内容
	createdGroups   map[string]groupInfo     // groupID -> groupInfo
	queriedUsers    []string                 // 查询过的用户
	shouldFail      map[string]error         // 方法名 -> 返回错误，空则不报错
}

type privateMsgCall struct {
	SenderID string
	TargetID string
	Content  string
}

type systemMsgCall struct {
	SenderID string
	TargetID string
	Content  string
}

type groupInfo struct {
	Name    string
	OwnerID string
}

func newMockIMClient() *mockIMClient {
	return &mockIMClient{
		registeredUsers: make(map[string]string),
		createdGroups:   make(map[string]groupInfo),
		shouldFail:      make(map[string]error),
	}
}

// mockFail 设置某个方法返回错误
func (m *mockIMClient) mockFail(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail[method] = err
}

func (m *mockIMClient) RegisterUser(userID, nickname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["RegisterUser"]; ok {
		return err
	}
	m.registeredUsers[userID] = nickname
	return nil
}

func (m *mockIMClient) SendPrivateMsg(senderID, targetID, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["SendPrivateMsg"]; ok {
		return err
	}
	m.privateMsgs = append(m.privateMsgs, privateMsgCall{
		SenderID: senderID,
		TargetID: targetID,
		Content:  content,
	})
	return nil
}

func (m *mockIMClient) SendSystemMsg(senderID, targetID, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["SendSystemMsg"]; ok {
		return err
	}
	m.systemMsgs = append(m.systemMsgs, systemMsgCall{
		SenderID: senderID,
		TargetID: targetID,
		Content:  content,
	})
	return nil
}

func (m *mockIMClient) SendGroupMsg(senderID, groupID, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["SendGroupMsg"]; ok {
		return err
	}
	m.groupMsgs = append(m.groupMsgs, content)
	return nil
}

func (m *mockIMClient) QueryUserInfo(userID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["QueryUserInfo"]; ok {
		return nil, err
	}
	m.queriedUsers = append(m.queriedUsers, userID)
	return map[string]interface{}{
		"user_id":  userID,
		"nickname": "mock_user_" + userID,
	}, nil
}

func (m *mockIMClient) CreateGroup(groupID, groupName, ownerID string, memberIDs []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["CreateGroup"]; ok {
		return err
	}
	m.createdGroups[groupID] = groupInfo{
		Name:    groupName,
		OwnerID: ownerID,
	}
	return nil
}

func (m *mockIMClient) QueryOnlineStatus(userIDs []string) (map[string]bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["QueryOnlineStatus"]; ok {
		return nil, err
	}
	result := make(map[string]bool, len(userIDs))
	for _, id := range userIDs {
		result[id] = true
	}
	return result, nil
}

func (m *mockIMClient) SendBroadcastMsg(senderID, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["SendBroadcastMsg"]; ok {
		return err
	}
	return nil
}

func (m *mockIMClient) AddBot(botID, nickname, webhookURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.shouldFail["AddBot"]; ok {
		return err
	}
	return nil
}

// assertRegisteredUser 断言某个用户已在 IM 注册
func (m *mockIMClient) assertRegisteredUser(userID, nickname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	got, ok := m.registeredUsers[userID]
	if !ok {
		return fmt.Errorf("用户 %s 未注册到 IM", userID)
	}
	if got != nickname {
		return fmt.Errorf("用户 %s 昵称不匹配: 期望 %s, 实际 %s", userID, nickname, got)
	}
	return nil
}

// assertPrivateMsgSent 断言某条私信已通过 IM 发送
func (m *mockIMClient) assertPrivateMsgSent(senderID, targetID, wantContent string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, msg := range m.privateMsgs {
		if msg.SenderID == senderID && msg.TargetID == targetID && msg.Content == wantContent {
			return nil
		}
	}
	return fmt.Errorf("未找到私信: sender=%s target=%s content=%s", senderID, targetID, wantContent)
}

// assertSystemMsgSent 断言某条系统消息已发送
func (m *mockIMClient) assertSystemMsgSent(targetID, wantContentSub string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, msg := range m.systemMsgs {
		if msg.TargetID == targetID && contains(msg.Content, wantContentSub) {
			return nil
		}
	}
	return fmt.Errorf("未找到系统消息: target=%s (包含: %s)", targetID, wantContentSub)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
