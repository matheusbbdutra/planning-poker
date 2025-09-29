package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/http/jsoncodec"
	"planning-poker/internal/infrastructure/persistence"
	"planning-poker/internal/ports/repository"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)


type RoomHandler struct {
	RedisRepository repository.RedisRepository
	Hub             Hub
}

type NewRoomRequest struct {
	UserName    string `json:"userName"`
	SessionName string `json:"sessionName"`
}

type JoinRoomRequest struct {
	UserName string `json:"userName"`
}

func NewRoomHandler(redisClient *redis.Client, hub Hub) *RoomHandler {
	return &RoomHandler{
		RedisRepository: persistence.NewRedisRepositoryImpl(redisClient),
		Hub:             hub,
	}
}

// NewRoom cria uma nova sala e retorna os detalhes da sala e do participante criado.
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

// JoinRoom adiciona um participante a uma sala existente.
func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	roomId := chi.URLParam(r, "roomId")
	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}
	log.Printf("Joining room %s with user %s", roomId, req.UserName)

	room, err := h.RedisRepository.GetRoom(r.Context(), roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return err
	}

	participant := entities.NewParticipant(req.UserName, false)
	room.AddParticipant(*participant)

	err = h.RedisRepository.UpdateRoom(r.Context(), *room)
	if err != nil {
		http.Error(w, "Failed to join room", http.StatusInternalServerError)
		return err
	}

	h.Hub.BroadcastToRoom(BroadcastMessage{
		RoomID: room.ID,
		Message: Message{
			Type:    "ROOM_STATE_UPDATED",
			Payload: jsoncodec.MustMarshal(room),
		},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(participant)
}
