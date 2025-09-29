package action

import (
	"context"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
)


type CardsAction struct {
	RedisClient persistence.RedisRepositoryImpl
}

func NewCardsAction(redisClient persistence.RedisRepositoryImpl) *CardsAction {
	return &CardsAction{
		RedisClient: redisClient,
	}
}

func (a *CardsAction) NewCardsExecute(roomId string, numberOfCards interface{}) (*entities.Room, error) {
	roomState, err := a.RedisClient.GetRoom(context.Background(), roomId)
	if err != nil {
		return nil, err
	}

	roomState.NumberOfCards = numberOfCards

	if err := a.RedisClient.UpdateRoom(context.Background(), *roomState); err != nil {
		return nil, err
	}
	return roomState, nil
}