package action

import (
	"context"
	"fmt"
	"planning-poker/internal/application/command"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
)

type TaskAction struct {
	RedisClient persistence.RedisRepositoryImpl
}

func NewTaskAction(redisClient persistence.RedisRepositoryImpl) *TaskAction {
	return &TaskAction{
		RedisClient: redisClient,
	}
}

func (a *TaskAction) NewTask(roomId string, task *command.CreateTaskRequestCommand) (*entities.Room, error) {
	roomState, err := a.RedisClient.GetRoom(context.Background(), roomId)
	if err != nil {
		return nil, err
	}
	newTask := entities.NewTask(task.Title)
	roomState.AddTask(*newTask)

	if err := a.RedisClient.UpdateRoom(context.Background(), *roomState); err != nil {
		return nil, err
	}
	return roomState, nil
}

func (a *TaskAction) AlterTaskStatus(roomId string, taskId string, status string) (*entities.Room, error) {
	roomState, err := a.RedisClient.GetRoom(context.Background(), roomId)
	if err != nil {
		return nil, err
	}
	task := roomState.GetTaskByID(taskId)
	task.SetStatus(status)

	if err := a.RedisClient.UpdateRoom(context.Background(), *roomState); err != nil {
		return nil, err
	}
	return roomState, nil
}

func (a *TaskAction) AddVote(roomId string, vote *command.UserVoteRequestCommand) (*entities.Room, error) {
	roomState, err := a.RedisClient.GetRoom(context.Background(), roomId)
	if err != nil {
		return nil, err
	}
	task := roomState.GetTaskByID(vote.TaskID)
	voteValue := fmt.Sprintf("%v", vote.Vote)
	task.AddVote(vote.UserID, voteValue)

	if err := a.RedisClient.UpdateRoom(context.Background(), *roomState); err != nil {
		return nil, err
	}
	return roomState, nil
}
