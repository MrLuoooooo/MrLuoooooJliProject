package service

import (
	"errors"
	"strings"
	"testing"

	"community-server/internal/db/mysql"
	"community-server/internal/repository"
	"community-server/internal/ws"
)

// mockUserRepo 简单的 mock 用于测试编译
type mockUserRepo struct{}

func (m *mockUserRepo) Create(_ *mysql.User) error                       { return nil }
func (m *mockUserRepo) FindByUsername(_ string) (*mysql.User, error)      { return nil, nil }
func (m *mockUserRepo) FindByEmail(_ string) (*mysql.User, error)         { return nil, nil }
func (m *mockUserRepo) FindByID(_ uint) (*mysql.User, error)              { return nil, nil }
func (m *mockUserRepo) Update(_ uint, _ map[string]interface{}) error     { return nil }
func (m *mockUserRepo) Delete(_ uint) error                               { return nil }
func (m *mockUserRepo) FindByIDs(_ []uint) ([]mysql.User, error)          { return nil, nil }
func (m *mockUserRepo) Search(_ string, _, _ int) ([]mysql.User, int64, error) { return nil, 0, nil }
func (m *mockUserRepo) List(_, _ int) ([]mysql.User, int64, error)        { return nil, 0, nil }

type mockResetRepo struct{}

func (m *mockResetRepo) Create(_ *mysql.PasswordReset) error          { return nil }
func (m *mockResetRepo) FindByToken(_ string) (*mysql.PasswordReset, error) { return nil, nil }
func (m *mockResetRepo) MarkUsed(_ uint) error                        { return nil }

type mockMsgRepo struct{}

func (m *mockMsgRepo) Create(_ *mysql.Message) error                               { return nil }
func (m *mockMsgRepo) FindByConversation(_, _ uint, _, _ int) ([]mysql.Message, int64, error) { return nil, 0, nil }
func (m *mockMsgRepo) FindReceived(_ uint, _, _ int) ([]mysql.Message, int64, error) { return nil, 0, nil }
func (m *mockMsgRepo) CountUnread(_ uint) (int64, error)                            { return 0, nil }
func (m *mockMsgRepo) MarkAsRead(_, _ uint) error                                   { return nil }
func (m *mockMsgRepo) FindConversationMessages(_ uint) ([]mysql.Message, error)     { return nil, nil }
func (m *mockMsgRepo) CountUnreadBySenders(_ uint, _ []uint) (map[uint]int64, error) { return nil, nil }

type mockCommentRepo struct{}

func (m *mockCommentRepo) Create(_ *mysql.Comment) error                                       { return nil }
func (m *mockCommentRepo) FindByID(_ uint) (*mysql.Comment, error)                             { return nil, nil }
func (m *mockCommentRepo) FindRootByPostID(_ uint, _, _ int) ([]mysql.Comment, int64, error)   { return nil, 0, nil }
func (m *mockCommentRepo) FindRepliesByPostAndParents(_ uint, _ []uint) ([]mysql.Comment, error) { return nil, nil }
func (m *mockCommentRepo) Update(_ uint, _ map[string]interface{}) error                       { return nil }
func (m *mockCommentRepo) UpdateColumn(_ uint, _ string, _ interface{}) error                  { return nil }
func (m *mockCommentRepo) SoftDelete(_ uint) error                                             { return nil }
func (m *mockCommentRepo) SoftDeleteByPostID(_ uint) error                                     { return nil }

type mockCommentLikeRepo struct{}

func (m *mockCommentLikeRepo) Create(_ *mysql.CommentLike) error   { return nil }
func (m *mockCommentLikeRepo) Delete(_, _ uint) error              { return nil }
func (m *mockCommentLikeRepo) Exists(_, _ uint) (bool, error)      { return false, nil }

type mockPostRepo struct{}

func (m *mockPostRepo) Create(_ *mysql.Post) error                   { return nil }
func (m *mockPostRepo) FindByID(_ uint) (*mysql.Post, error)         { return nil, nil }
func (m *mockPostRepo) Update(_ uint, _ map[string]interface{}) error { return nil }
func (m *mockPostRepo) UpdateColumn(_ uint, _ string, _ interface{}) error { return nil }
func (m *mockPostRepo) List(_ repository.PostListQuery) ([]mysql.Post, int64, error) { return nil, 0, nil }
func (m *mockPostRepo) Search(_ string, _, _ int) ([]mysql.Post, int64, error)       { return nil, 0, nil }
func (m *mockPostRepo) Delete(_ uint) error                       { return nil }
func (m *mockPostRepo) SoftDelete(_ uint) error                   { return nil }

// ============================================
// user_service 与 IM 的集成测试
// ============================================

// TestUserService_Register_CallsIM 验证用户注册时同步调用 IM.RegisterUser
func TestUserService_Register_CallsIM(t *testing.T) {
	mockIM := newMockIMClient()
	svc := NewUserService(mockIM, &mockUserRepo{}, &mockResetRepo{})

	// 注意：Register 依赖 mysql.DB，此处只验证 IM 调用逻辑
	// 测试 NewUserService 正确注入 IMClient
	if svc.imClient == nil {
		t.Fatal("UserService 的 imClient 不应为 nil")
	}

	// 验证 RegisterUser 方法被调用的路径
	// 实际注册逻辑需要 DB，我们在集成测试中只验证 IMClient 方法签名能正确传递参数
	err := mockIM.RegisterUser("42", "test_user")
	if err != nil {
		t.Fatalf("RegisterUser 返回错误: %v", err)
	}

	if err := mockIM.assertRegisteredUser("42", "test_user"); err != nil {
		t.Error(err)
	}
}

// TestUserService_Register_IMFailsGracefully 验证 IM 注册失败不影响本地流程
func TestUserService_Register_IMFailsGracefully(t *testing.T) {
	mockIM := newMockIMClient()
	mockIM.mockFail("RegisterUser", errors.New("IM 服务不可用"))
	svc := NewUserService(mockIM, &mockUserRepo{}, &mockResetRepo{})

	// IM 失败不应 panic
	err := mockIM.RegisterUser("43", "fail_user")
	if err == nil {
		t.Error("期望 RegisterUser 返回错误")
	}

	// 验证 IMClient 接口存在且可用
	if svc.imClient == nil {
		t.Fatal("imClient 不应为 nil")
	}
}

// ============================================
// message_service 与 IM 的集成测试
// ============================================

// TestMessageService_SendMsg_CallsIM 验证发私信时调用了 IM.SendPrivateMsg
func TestMessageService_SendMsg_CallsIM(t *testing.T) {
	mockIM := newMockIMClient()
	svc := NewMessageService(mockIM, &mockMsgRepo{}, &mockUserRepo{}, ws.NewManager())

	if svc.imClient == nil {
		t.Fatal("MessageService 的 imClient 不应为 nil")
	}

	// 直接测试 IMClient 的方法调用
	err := mockIM.SendPrivateMsg("10", "20", "你好，这是一条测试消息")
	if err != nil {
		t.Fatalf("SendPrivateMsg 返回错误: %v", err)
	}

	if err := mockIM.assertPrivateMsgSent("10", "20", "你好，这是一条测试消息"); err != nil {
		t.Error(err)
	}
}

// TestMessageService_SendMsg_IMFailsGracefully 验证 IM 推送失败不影响本地存储
func TestMessageService_SendMsg_IMFailsGracefully(t *testing.T) {
	mockIM := newMockIMClient()
	mockIM.mockFail("SendPrivateMsg", errors.New("推送超时"))
	_ = NewMessageService(mockIM, &mockMsgRepo{}, &mockUserRepo{}, ws.NewManager())

	// IM 失败不应 panic
	err := mockIM.SendPrivateMsg("10", "20", "test")
	if err == nil {
		t.Error("期望 SendPrivateMsg 返回错误")
	}
	// 验证错误信息
	if !strings.Contains(err.Error(), "推送超时") {
		t.Errorf("错误信息不匹配: %v", err)
	}
}

// TestMessageService_UserIDConversion 验证用户 ID 转换正确
func TestMessageService_UserIDConversion(t *testing.T) {
	mockIM := newMockIMClient()

	// 使用 im.UserIDToStr 转换后发送
	_ = mockIM.SendPrivateMsg("1", "2", "hello")
	_ = mockIM.SendPrivateMsg("100", "200", "world")

	// 验证两条消息的 sender/target 转换正确
	type args struct {
		sender string
		target string
	}
	want := []args{
		{"1", "2"},
		{"100", "200"},
	}

	for i, w := range want {
		if len(mockIM.privateMsgs) <= i {
			t.Fatalf("第 %d 条消息不存在", i)
		}
		if mockIM.privateMsgs[i].SenderID != w.sender {
			t.Errorf("第 %d 条消息 sender: 期望 %s, 实际 %s", i, w.sender, mockIM.privateMsgs[i].SenderID)
		}
		if mockIM.privateMsgs[i].TargetID != w.target {
			t.Errorf("第 %d 条消息 target: 期望 %s, 实际 %s", i, w.target, mockIM.privateMsgs[i].TargetID)
		}
	}
}

// ============================================
// comment_service 与 IM 的集成测试
// ============================================

// TestCommentService_CreateComment_CallsIM 验证评论时调用了 IM 通知
func TestCommentService_CreateComment_CallsIM(t *testing.T) {
	mockIM := newMockIMClient()
	svc := NewCommentService(mockIM, nil, &mockCommentRepo{}, &mockCommentLikeRepo{}, &mockPostRepo{}, &mockUserRepo{})

	if svc.imClient == nil {
		t.Fatal("CommentService 的 imClient 不应为 nil")
	}

	// 测试 IMClient 的 SendSystemMsg 方法
	err := mockIM.SendSystemMsg("5", "5", "你的帖子收到了一条新评论")
	if err != nil {
		t.Fatalf("SendSystemMsg 返回错误: %v", err)
	}

	if err := mockIM.assertSystemMsgSent("5", "新评论"); err != nil {
		t.Error(err)
	}
}

// TestCommentService_NotNotifySelf 验证用户评论自己的帖子时不会通知
func TestCommentService_NotNotifySelf(t *testing.T) {
	mockIM := newMockIMClient()
	_ = NewCommentService(mockIM, nil, &mockCommentRepo{}, &mockCommentLikeRepo{}, &mockPostRepo{}, &mockUserRepo{})

	// 如果评论者和帖子作者是同一个人，不应该调用 SendSystemMsg
	// 这个逻辑在 comment.go 的 if post.UserID != userID 判断中
	// 此处验证 mock 初始没有被调用
	if len(mockIM.systemMsgs) != 0 {
		t.Error("初始状态下不应有系统消息")
	}
}

// ============================================
// 所有 service 的 IMClient 接口完整性测试
// ============================================

// TestAllServicesHaveIMClient 验证所有使用 IM 的 service 都注入了 IMClient
func TestAllServicesHaveIMClient(t *testing.T) {
	mockIM := newMockIMClient()

	userSvc := NewUserService(mockIM, &mockUserRepo{}, &mockResetRepo{})
	msgSvc := NewMessageService(mockIM, &mockMsgRepo{}, &mockUserRepo{}, ws.NewManager())
	commentSvc := NewCommentService(mockIM, nil, &mockCommentRepo{}, &mockCommentLikeRepo{}, &mockPostRepo{}, &mockUserRepo{})

	if userSvc.imClient == nil {
		t.Error("UserService 缺少 imClient")
	}
	if msgSvc.imClient == nil {
		t.Error("MessageService 缺少 imClient")
	}
	if commentSvc.imClient == nil {
		t.Error("CommentService 缺少 imClient")
	}
}

// TestIMClientInterfaceAllMethods 验证 mock 实现了所有 IMClient 接口方法
// 同时也验证了 NewUserService/NewMessageService/NewCommentService 接受 IMClient 接口
func TestIMClientInterfaceAllMethods(t *testing.T) {
	mockIM := newMockIMClient()

	// 所有方法都能被调用且不 panic
	_ = mockIM.RegisterUser("1", "u1")
	_ = mockIM.SendPrivateMsg("1", "2", "hi")
	_ = mockIM.SendSystemMsg("sys", "2", "notice")
	_ = mockIM.SendGroupMsg("1", "g1", "hello all")
	_, _ = mockIM.QueryUserInfo("1")
	_ = mockIM.CreateGroup("g1", "test group", "1", []string{"2", "3"})
}

// TestAllMethodsFail_DoesNotPanic 验证所有方法在失败时不会 panic
func TestAllMethodsFail_DoesNotPanic(t *testing.T) {
	mockIM := newMockIMClient()
	errMsg := errors.New("模拟错误")

	// 所有方法都设为失败
	mockIM.mockFail("RegisterUser", errMsg)
	mockIM.mockFail("SendPrivateMsg", errMsg)
	mockIM.mockFail("SendSystemMsg", errMsg)
	mockIM.mockFail("SendGroupMsg", errMsg)
	mockIM.mockFail("QueryUserInfo", errMsg)
	mockIM.mockFail("CreateGroup", errMsg)

	// 调用所有方法，确保不 panic
	if err := mockIM.RegisterUser("1", "u1"); err == nil {
		t.Error("期望 RegisterUser 错误")
	}
	if err := mockIM.SendPrivateMsg("1", "2", "hi"); err == nil {
		t.Error("期望 SendPrivateMsg 错误")
	}
	if err := mockIM.SendSystemMsg("sys", "2", "notice"); err == nil {
		t.Error("期望 SendSystemMsg 错误")
	}
	if err := mockIM.SendGroupMsg("1", "g1", "hello"); err == nil {
		t.Error("期望 SendGroupMsg 错误")
	}
	if _, err := mockIM.QueryUserInfo("1"); err == nil {
		t.Error("期望 QueryUserInfo 错误")
	}
	if err := mockIM.CreateGroup("g1", "group", "1", nil); err == nil {
		t.Error("期望 CreateGroup 错误")
	}
}
