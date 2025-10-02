package action

import (
	"context"
	"planning-poker/internal/application/command"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
)


type RoomAction struct {
	RedisClient persistence.RedisRepositoryImpl
}

func NewRoomAction(redisClient persistence.RedisRepositoryImpl) *RoomAction {
	return &RoomAction{
		RedisClient: redisClient,
	}
}

func (a *RoomAction) ParticipantJoinRoom(roomId string, participantJoin *command.ParticipantJoinRoomCommand) (*entities.Room, error) {
	roomState, err := a.RedisClient.GetRoom(context.Background(), roomId)
	if err != nil {
		return nil, err
	}

	participant := entities.NewParticipant(participantJoin.UserName, false)
	roomState.AddParticipant(*participant)

	if err := a.RedisClient.UpdateRoom(context.Background(), *roomState); err != nil {
		return nil, err
	}
	return roomState, nil
}