package entities



type Participant struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IsScrumMaster bool   `json:"isScrumMaster"`
	IsConnected   bool   `json:"isConnected"`
}

func NewParticipant(name string, isScrumMaster bool) *Participant {
	id, err := NewID()
	if err != nil {
		panic(err)
	}
	return &Participant{
		ID:            id,
		Name:          name,
		IsScrumMaster: isScrumMaster,
		IsConnected:   true,
	}
}
