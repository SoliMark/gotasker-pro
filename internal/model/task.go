package model

import "time"

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Title     string `gorm:"not null"`
	Content   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

const (
	TaskStatusPending = "pending"
	TaskStatusDone    = "done"
)
