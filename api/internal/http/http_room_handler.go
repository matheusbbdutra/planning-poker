package http_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
	"planning-poker/internal/ports/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type HttpRoomHandler struct {
	RedisRepository repository.RedisRepository
}

type NewRoomRequest struct {
	UserName    string `json:"user_name"`
	SessionName string `json:"session_name"`
}

func NewRoomHandler(redisClient *redis.Client) *HttpRoomHandler {
	return &HttpRoomHandler{
		RedisRepository: persistence.NewRedisRepositoryImpl(redisClient),
	}
}

func (h *HttpRoomHandler) NewRoom(w http.ResponseWriter, r *http.Request) error {
	var req NewRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}
	log.Printf("%v", req)

	room := entities.Room{
		ID:   uuid.New(),
		Name: req.SessionName,
		Participants: []entities.Participant{
			*entities.NewParticipant(req.UserName, true),
		},
	}

	err := h.RedisRepository.CreateRoom(r.Context(), room)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{"room": room}
	return json.NewEncoder(w).Encode(response)
}
