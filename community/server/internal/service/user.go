package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"community-server/DB/mysql"
	"community-server/internal/model"
	"community-server/pkg/jwt"

	"go.uber.org/zap"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *model.RegisterRequest) (uint, error) {
	var existingUser mysql.User
	result := mysql.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return 0, errors.New("用户名已存在")
	}

	if req.Email != "" {
		result = mysql.DB.Where("email = ?", req.Email).First(&existingUser)
		if result.Error == nil {
			return 0, errors.New("邮箱已被注册")
		}
	}

	hashedPassword := hashPassword(req.Password)

	adminType := 0
	if req.AdminType == 1 {
		adminType = 1
	}

	user := mysql.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Nickname:  req.Nickname,
		Avatar:    "",
		Bio:       "",
		AdminType: adminType,
		Status:    1,
	}

	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	result = mysql.DB.Create(&user)
	if result.Error != nil {
		zap.S().Error("用户注册失败", "username", req.Username, "error", result.Error)
		return 0, errors.New("注册失败")
	}

	zap.S().Info("用户注册成功", "userId", user.ID, "username", req.Username, "adminType", adminType)
	return user.ID, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	var user mysql.User
	result := mysql.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Password != hashPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		zap.S().Error("生成令牌失败", "userId", user.ID, "error", err)
		return nil, errors.New("生成令牌失败")
	}

	now := time.Now()
	mysql.DB.Model(&user).Update("last_login", &now)

	zap.S().Info("用户登录成功", "userId", user.ID, "username", user.Username)

	return &model.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *UserService) GetUserByID(userID uint) (*mysql.User, error) {
	var user mysql.User
	result := mysql.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uint, req *model.UpdateProfileRequest) error {
	var user mysql.User
	result := mysql.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
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
		mysql.DB.Model(&user).Updates(updates)
	}

	return nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
