package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"community-server/internal/db/mysql"
	"community-server/internal/im"
	"community-server/internal/model"
	"community-server/internal/repository"
	"community-server/pkg/jwt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	imClient  im.IMClient
	userRepo  repository.UserRepository
	resetRepo repository.PasswordResetRepository
}

func NewUserService(imClient im.IMClient, userRepo repository.UserRepository, resetRepo repository.PasswordResetRepository) *UserService {
	return &UserService{imClient: imClient, userRepo: userRepo, resetRepo: resetRepo}
}

func (s *UserService) Register(req *model.RegisterRequest) (uint, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return 0, errors.New("用户名已存在")
	}
	if req.Email != "" {
		if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
			return 0, errors.New("邮箱已被注册")
		}
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return 0, errors.New("密码加密失败")
	}

	user := mysql.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		Avatar:    "",
		Bio:       "",
		AdminType: 0,
		Status:    1,
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.userRepo.Create(&user); err != nil {
		zap.S().Error("用户注册失败", "username", req.Username, "error", err)
		return 0, errors.New("注册失败")
	}

	if err := s.imClient.RegisterUser(im.UserIDToStr(user.ID), user.Nickname); err != nil {
		zap.S().Warn("IM注册失败（不影响本地注册）", "userId", user.ID, "error", err)
	}

	zap.S().Info("用户注册成功", "userId", user.ID, "username", req.Username)
	return user.ID, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		zap.S().Error("生成令牌失败", "userId", user.ID, "error", err)
		return nil, errors.New("生成令牌失败")
	}

	now := time.Now()
	s.userRepo.Update(user.ID, map[string]interface{}{"last_login": &now})

	zap.S().Info("用户登录成功", "userId", user.ID, "username", user.Username)
	return &model.LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		AdminType: user.AdminType,
	}, nil
}

func (s *UserService) GetUserByID(userID uint) (*model.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &model.UserProfileResponse{
		ID: user.ID, Username: user.Username, Nickname: user.Nickname,
		Avatar: user.Avatar, Bio: user.Bio, Email: user.Email, Status: user.Status, AdminType: user.AdminType,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, req *model.UpdateProfileRequest) error {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return errors.New("用户不存在")
	}
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if len(updates) > 0 {
		return s.userRepo.Update(userID, updates)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ForgotPassword 生成密码重置令牌。开发环境令牌写入日志；生产环境应通过邮件发送。
func (s *UserService) ForgotPassword(email string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("该邮箱未注册")
	}

	token := generateResetToken()
	reset := mysql.PasswordReset{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
	}
	if err := s.resetRepo.Create(&reset); err != nil {
		zap.S().Error("创建密码重置令牌失败", "email", email, "error", err)
		return errors.New("系统繁忙，请稍后重试")
	}
	zap.S().Info("密码重置令牌", "email", email, "token", token)
	return nil
}

// ResetPassword 校验令牌并重置密码
func (s *UserService) ResetPassword(token, newPassword string) error {
	reset, err := s.resetRepo.FindByToken(token)
	if err != nil {
		return errors.New("无效的重置令牌")
	}
	if reset.Used {
		return errors.New("重置令牌已使用")
	}
	if time.Now().Unix() > reset.ExpiresAt {
		return errors.New("重置令牌已过期")
	}

	hashed, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user, err := s.userRepo.FindByEmail(reset.Email)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := s.userRepo.Update(user.ID, map[string]interface{}{"password": hashed}); err != nil {
		zap.S().Error("更新密码失败", "userId", user.ID, "error", err)
		return errors.New("密码重置失败")
	}
	if err := s.resetRepo.MarkUsed(reset.ID); err != nil {
		zap.S().Warn("标记令牌已使用失败", "resetID", reset.ID, "error", err)
	}

	zap.S().Info("密码重置成功", "userId", user.ID)
	return nil
}

func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
