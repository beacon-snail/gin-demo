package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Age       int            `json:"age"`
	Gender    string         `json:"gender" gorm:"size:10;default:'unknown'"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Status    int            `json:"status" gorm:"default:1"` // 1: active, 0: inactive
}

func (User) TableName() string {
	return "users"
}
