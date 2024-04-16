package user

import "github.com/Conty111/SuperCalculator/back-end/models"

type MsgResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

type UserInfo struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
}

// UserResponse is a declaration for a common user response
type UserResponse struct {
	Status string `json:"status"`
	UserInfo
}

type UsersListResponse struct {
	Status string     `json:"status"`
	Users  []UserInfo `json:"users"`
}

// AuthResponse is a declaration for an auth response
type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	UserInfo
}
