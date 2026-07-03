package dto

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Age      int    `json:"age" binding:"omitempty,min=0,max=150"`
	Gender   string `json:"gender" binding:"omitempty,oneof=male female unknown"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Avatar   string `json:"avatar"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Nickname  string     `json:"nickname"`
	Age       int        `json:"age"`
	Gender    string     `json:"gender"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

type UserListResponse struct {
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Data  []UserResponse `json:"data"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
