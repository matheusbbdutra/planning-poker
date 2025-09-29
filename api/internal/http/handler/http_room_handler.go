package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
	"planning-poker/internal/ports/repository"

	"github.com/redis/go-redis/v9"
)

type RoomHandler struct {
	RedisRepository repository.RedisRepository
}

type NewRoomRequest struct {
	UserName    string `json:"userName"`
	SessionName string `json:"sessionName"`
}

func NewRoomHandler(redisClient *redis.Client) *RoomHandler {
	return &RoomHandler{
		RedisRepository: persistence.NewRedisRepositoryImpl(redisClient),
	}
}

func (h *RoomHandler) NewRoom(w http.ResponseWriter, r *http.Request) error {
	var req NewRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}
	log.Printf("%v", req)

	room := entities.NewRoom(req.SessionName)
	participant := entities.NewParticipant(req.UserName, true)
	room.AddParticipant(*participant)

	err := h.RedisRepository.CreateRoom(r.Context(), *room)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{"room": room, "participant": participant} // O frontend espera o objeto room
	return json.NewEncoder(w).Encode(response)
}
