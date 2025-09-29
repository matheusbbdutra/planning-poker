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

	key := "room:" + room.ID
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
		return nil, err
	}

	log.Printf("GetRoom found key=%s value=%s", key, roomData)

	var room entities.Room
	if err := json.Unmarshal([]byte(roomData), &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RedisRepositoryImpl) UpdateRoom(ctx context.Context, room entities.Room) error {
	roomData, err := json.Marshal(room)
	if err != nil {
		return err
	}

	key := "room:" + room.ID
	log.Printf("UpdateRoom set key=%s value=%s", key, string(roomData))

	if err := r.Redis.Set(ctx, key, roomData, time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepositoryImpl) RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error {
    room, err := r.GetRoom(ctx, roomId)
    if err != nil {
        return err
    }

    // Remove participant
    updatedParticipants := []entities.Participant{}
    for _, participant := range room.Participants {
        if participant.ID != userId {
            updatedParticipants = append(updatedParticipants, participant)
        }
    }
    room.Participants = updatedParticipants

    // Update room in Redis
    return r.UpdateRoom(ctx, *room)
}

func (r *RedisRepositoryImpl) ConnectedUser(ctx context.Context, roomId string, userId string, isConnected bool) error {
    room, err := r.GetRoom(ctx, roomId)
    if err != nil {
        return err
    }   
    // Mark participant as disconnected
    for i, participant := range room.Participants {
        if participant.ID == userId {
            room.Participants[i].IsConnected = isConnected
            break
        }
    }   
    // Update room in Redis
    return r.UpdateRoom(ctx, *room)
}