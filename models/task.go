package models

import (
	"time"
)

// Task "Object
type Task struct {
	ID        uint      `gorm:primaryKey json:"id"`
	Title     string    `json:"title" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Completed bool      `json:"completed"`
}
