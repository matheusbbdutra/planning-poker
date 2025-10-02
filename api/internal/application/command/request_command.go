package command

type CreateTaskRequestCommand struct {
	Title string `json:"title"`
}

type UserVoteRequestCommand struct {
	UserID string      `json:"userId"`
	TaskID string      `json:"taskId"`
	Vote   interface{} `json:"vote"`
}

type SetNumberOfCardsRequestCommand struct {
	NumberOfCards interface{} `json:"numberOfCards"`
}

type AlterTaskStatusRequestCommand struct {
	TaskID       string `json:"taskId"`
	VotingStatus string `json:"votingStatus"`
}

type ParticipantJoinRoomCommand struct {
	UserName string `json:"userName"`
}