package models

import (
	"time"
)

// Task "Object
type Task struct {
	ID        uint      `gorm:primaryKey json:"id"`
	Title     string    `json:"title" binding:"required"`
	UserId    string    `gorm: "index:idx_user_id" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Completed bool      `json:"completed"`
}
