package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // tidak pernah dikirim ke client
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserCreateInput struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"omitempty,oneof=owner worker"`
}

type UserUpdateInput struct {
	Name     string  `json:"name" binding:"required,max=100"`
	Email    string  `json:"email" binding:"required,email,max=150"`
	Role     string  `json:"role" binding:"required,oneof=owner worker"`
	Password *string `json:"password" binding:"omitempty,min=8,max=72"` // diisi hanya kalau mau reset
}

type UserStatusInput struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type UserFilter struct {
	Search   string
	Role     string
	IsActive *bool
	Page     int
	Limit    int
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
