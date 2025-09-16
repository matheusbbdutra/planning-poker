package persistence

import (
	"context"
	"planning-poker/internal/domain/entities"

	"github.com/redis/go-redis/v9"
)

type RedisRepositoryImpl struct {
	Redis *redis.Client
}

func NewRedisRepositoryImpl(redisClient *redis.Client) *RedisRepositoryImpl {
	return &RedisRepositoryImpl{
		Redis: redisClient,
	}
}

func (r *RedisRepositoryImpl) CreateRoom(ctx context.Context, room entities.Room) error {
	if err := r.Redis.Set(ctx, room.ID.String(), room.Name, 3600).Err(); err != nil {	
		return err
	}
	return nil
}
