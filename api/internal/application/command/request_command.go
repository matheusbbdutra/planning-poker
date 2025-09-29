package command


type CreateTaskRequestCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateParticipantRequestCommand struct {
	Name string `json:"name"`
}

type UserVoteRequestCommand struct {
	ParticipantID string `json:"participantId"`
	TaskID        string `json:"taskId"`
	CardValue    int    `json:"cardValue"`
}

type SetNumberOfCardsRequestCommand struct {
	NumberOfCards interface{} `json:"numberOfCards"`
}

