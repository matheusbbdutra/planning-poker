package repository

import (
	"context"
	"planning-poker/internal/domain/entities"
)

type RedisRepository interface {
	CreateRoom(ctx context.Context, room entities.Room) error
	GetRoom(ctx context.Context, roomId string) (*entities.Room, error)
	UpdateRoom(ctx context.Context, room entities.Room) error
	RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error
	ConnectedUser(ctx context.Context, roomId string, userId string, isConnected bool) error
}