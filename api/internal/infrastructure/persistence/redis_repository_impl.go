package persistence

import (
	"context"
	"encoding/json"
	"log"
	"planning-poker/internal/domain/entities"
	"time"

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
    roomData, err := json.Marshal(room)
    if err != nil {
        return err
    }

    key := "room:" + room.ID.String()
    log.Printf("CreateRoom set key=%s value=%s", key, string(roomData))

    if err := r.Redis.Set(ctx, key, roomData, time.Hour).Err(); err != nil {
        return err
    }
    return nil
}

func (r *RedisRepositoryImpl) GetRoom(ctx context.Context, roomId string) (*entities.Room, error) {
    key := "room:" + roomId
    roomData, err := r.Redis.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            log.Printf("chave %s não encontrada, tentando sem prefixo...", key)
            roomData2, err2 := r.Redis.Get(ctx, roomId).Result()
            if err2 == nil {
                roomData = roomData2
                key = roomId
            } else if err2 == redis.Nil {
                log.Printf("room %s não encontrada (tentadas: %s e %s)", roomId, "room:"+roomId, roomId)
                return nil, redis.Nil
            } else {
                return nil, err2
            }
        } else {
            return nil, err
        }
    }

    log.Printf("GetRoom found key=%s value=%s", key, roomData)

    var room entities.Room
    if err := json.Unmarshal([]byte(roomData), &room); err != nil {
        return nil, err
    }
    return &room, nil
}