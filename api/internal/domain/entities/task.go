package entities

import (
	"log"
	"planning-poker/internal/domain/enum"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	VoteCounts  map[string]int `json:"voteCounts"`
	IsCompleted bool `json:"isCompleted"`
	Votes       map[string]string `json:"votes"`
	VotingStatus string `json:"votingStatus"`
}

func NewTask(title string) *Task {
	id, err := NewID()
	if err != nil {
		log.Printf("Error generating task ID: %v", err)
		panic(err)
	}
	return &Task{
		ID:          id,
		Title:       title,
		VoteCounts:  make(map[string]int),
		IsCompleted: false,
		Votes:       make(map[string]string),
		VotingStatus: enum.STATUS_PENDING,
	}
}

func (t *Task) AddVote(participantID string, cardValue string) {
	if t.Votes == nil {
		t.Votes = make(map[string]string)
	}
	if t.VoteCounts == nil {
		t.VoteCounts = make(map[string]int)
	}

	if previousValue, alreadyVoted := t.Votes[participantID]; alreadyVoted {
		if previousValue == cardValue {
			return
		}

		if count, ok := t.VoteCounts[previousValue]; ok {
			if count <= 1 {
				delete(t.VoteCounts, previousValue)
			} else {
				t.VoteCounts[previousValue] = count - 1
			}
		}
	}

	t.Votes[participantID] = cardValue
	t.VoteCounts[cardValue] = t.VoteCounts[cardValue] + 1
}

func (t *Task) SetStatus(status string) {
	t.VotingStatus = status

	if status == enum.STATUS_COMPLETED {
		t.IsCompleted = true
	}
}
