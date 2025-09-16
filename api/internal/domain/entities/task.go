package entities

import "github.com/google/uuid"

const (
	STATUS_PENDING   = "pending"
	STATUS_COMPLETED = "completed"
)

type Task struct {
	ID          string
	Title       string
	VoteCounts  map[string]int
	isCompleted bool
	Votes       map[uuid.UUID]int16
	Status      string
}

func NewTask(title string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		VoteCounts:  make(map[string]int),
		isCompleted: false,
		Votes:      make(map[uuid.UUID]int16),
		Status:     STATUS_PENDING,
	}
}


