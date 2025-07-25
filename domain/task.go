package domain

import (
	"time"
)

const (
	StatusPending    = "Pending"
	StatusInProgress = "In Progress"
	StatusCompleted  = "Completed"
)

type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}
