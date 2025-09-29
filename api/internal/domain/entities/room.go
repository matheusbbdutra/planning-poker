package entities


type Room struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	Participants          []Participant `json:"participants"`
	Tasks                 []Task        `json:"tasks"`
	NumberOfTaskCompleted int           `json:"numberOfTaskCompleted"`
	NumberOfCards         interface{}   `json:"numberOfCards"`
}

func NewRoom(name string) *Room {
	id, err := NewID()
	if err != nil {
		panic(err)
	}
	return &Room{
		ID:                    id,
		Name:                  name,
		Participants:          []Participant{},
		Tasks:                 []Task{},
		NumberOfTaskCompleted: 0,
		NumberOfCards:         interface{}(nil),
	}
}

func (r *Room) AddParticipant(participant Participant) {
	r.Participants = append(r.Participants, participant)
}

func (r *Room) AddTask(task Task) {
	r.Tasks = append(r.Tasks, task)
}

func (r *Room) GetTaskByID(taskID string) *Task {
	for i := range r.Tasks {
		if r.Tasks[i].ID == taskID {
			return &r.Tasks[i]
		}
	}
	return nil
}

func (r *Room) GetParticipantByID(participantID string) *Participant {
	for i := range r.Participants {
		if r.Participants[i].ID == participantID {
			return &r.Participants[i]
		}
	}
	return nil
}

func (r *Room) SetNumberOfCards(number []int) {
	r.NumberOfCards = number
}
