package model

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"email"`
	Nickname string `json:"nickname" binding:"max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Bio      string `json:"bio" binding:"max=500"`
}
