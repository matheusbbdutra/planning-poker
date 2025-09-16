package entities

import "github.com/google/uuid"

type Participant struct {
	ID   uuid.UUID
	Name string
	IsScrumMaster bool 
}

func NewParticipant(name string, isScrumMaster bool) *Participant {
	return &Participant{
		ID:   uuid.New(),
		Name: name,
		IsScrumMaster: isScrumMaster,
	}
}
