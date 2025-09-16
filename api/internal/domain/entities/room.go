package entities

import "github.com/google/uuid"

type Room struct {
	ID          uuid.UUID
	Name        string
	Participants []Participant
	Tasks       []Task
	NumberOfTaskCompleted int
	NumberOfCards []int
}

func NewRoom(name string) *Room {
	return &Room{
		ID:                   uuid.New(),
		Name:                 name,
		Participants:         []Participant{},
		Tasks:               []Task{},
		NumberOfTaskCompleted: 0,
		NumberOfCards:       []int{},
	}
}

func (r *Room) AddParticipant(participant Participant) {
	r.Participants = append(r.Participants, participant)
}

func (r *Room) AddTask(task Task) {
	r.Tasks = append(r.Tasks, task)
}

func (r *Room) GetTaskByID(taskID string) *Task {
	for _, task := range r.Tasks {
		if task.ID == taskID {
			return &task
		}
	}
	return nil
}

func (r *Room) GetParticipantByID(participantID uuid.UUID) *Participant {
	for _, participant := range r.Participants {
		if participant.ID == participantID {
			return &participant
		}
	}
	return nil
}

func (r *Room) CompleteTask(taskID string) {
	task := r.GetTaskByID(taskID)
	if task != nil {
		task.Status = "completed"
		r.NumberOfTaskCompleted++
	}
}

func (r *Room) SetNumberOfCards(number []int) {
	r.NumberOfCards = number
}
