package entities

import "github.com/google/uuid"

type Participant struct {
	ID   uuid.UUID `json:"id"`
	Name string     `json:"name"`
	IsScrumMaster bool       `json:"isScrumMaster"`
}

func NewParticipant(name string, isScrumMaster bool) *Participant {
	return &Participant{
		ID:   uuid.New(),
		Name: name,
		IsScrumMaster: isScrumMaster,
	}
}
