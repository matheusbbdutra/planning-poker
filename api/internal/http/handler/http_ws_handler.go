package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"planning-poker/internal/application/command"
	handleraction "planning-poker/internal/application/handler_action"
	"planning-poker/internal/domain/entities"
	"planning-poker/internal/infrastructure/persistence"
	"planning-poker/internal/ports/repository"
	"planning-poker/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permite todas as origens por enquanto. Em produção, restrinja!
		return true
	},
}

type Hub interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastToRoom(message BroadcastMessage)
}


type Client struct {
	Conn   *websocket.Conn
	Hub    	Hub
	RoomID string
	UserID string
	Redis  repository.RedisRepository
}

type BroadcastMessage struct {
	RoomID  string
	Message Message
}

type HttpWsHandler struct {
	Hub             Hub
	RedisRepository repository.RedisRepository
	CardsAction *handleraction.CardsAction
	TaskAction           *handleraction.TaskAction
}

type Message struct {
	Type    string      `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewWsHandler(hub Hub, redisClient *redis.Client) *HttpWsHandler {
	redis := persistence.NewRedisRepositoryImpl(redisClient)

	return &HttpWsHandler{
		Hub:             hub,
		RedisRepository: redis,
		CardsAction:     handleraction.NewCardsAction(*redis),
		TaskAction:      handleraction.NewTaskAction(*redis),
	}
}

// HandleWebSocket lida com a conexão WebSocket.
func (h *HttpWsHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomId := chi.URLParam(r, "roomId")
	userId := chi.URLParam(r, "userId")

	client := &Client{
		Conn:   conn,
		Hub:    h.Hub,
		RoomID: roomId,
		UserID: userId,
		Redis:  h.RedisRepository,
	}

	// Registrar o client no hub
	h.Hub.RegisterClient(client)

	// Iniciar a leitura de mensagens do WebSocket
	go func() {
		defer func() {
			h.Hub.UnregisterClient(client)
			conn.Close()
		}()
		for {
			_, payload, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var msg Message
			if err := json.Unmarshal(payload, &msg); err != nil {
				log.Println("unmarshal:", err)
				break
			}

			log.Printf("Received message: Type: %s, Payload: %s", msg.Type, string(msg.Payload))

			switch msg.Type {
			case "UPDATE_CARDS":
				var cards command.SetNumberOfCardsRequestCommand
				if err := json.Unmarshal(msg.Payload, &cards); err != nil {
					log.Println("Invalid payload for UPDATE_CARDS:", err)
					continue
				}

				roomState, err := h.CardsAction.NewCardsExecute(roomId, cards.NumberOfCards)
				if err != nil {
					log.Println("Error updating cards:", err)
					continue
				}
				log.Printf("Updated room state: %+v", roomState)

				h.roomUpdate(roomId, *roomState)
			case "ADD_TASK":
				var task command.CreateTaskRequestCommand
				if err := json.Unmarshal(msg.Payload, &task); err != nil {
					log.Println("Invalid payload for ADD_TASK:", err)
					continue
				}

				roomState, err := h.TaskAction.NewTask(roomId, &task)
				if err != nil {
					log.Println("Error adding task:", err)
					continue
				}
				log.Printf("Updated room state: %+v", roomState)

				h.roomUpdate(roomId, *roomState)
				
			case "ON_VOTING":
				var status command.AlterTaskStatusRequestCommand
				if err := json.Unmarshal(msg.Payload, &status); err != nil {
					log.Println("Invalid payload for ON_VOTING:", err)
					continue
				}

				roomState, err := h.TaskAction.AlterTaskStatus(roomId, status.TaskID, status.VotingStatus)
				if err != nil {
					log.Println("Error changing task status:", err)
					continue
				}
				log.Printf("Updated room state: %+v", roomState)

				h.roomUpdate(roomId, *roomState)
			case "USER_VOTE":
				var vote command.UserVoteRequestCommand
				if err := json.Unmarshal(msg.Payload, &vote); err != nil {
					log.Println("Invalid payload for USER_VOTE:", err)
					continue
				}

				roomState, err := h.TaskAction.AddVote(roomId, &vote)
				if err != nil {
					log.Println("Error adding vote:", err)
					continue
				}
				log.Printf("Updated room state after vote: %+v", roomState)

				h.roomUpdate(roomId, *roomState)
			}
		}
	}()
}

func (h *HttpWsHandler) roomUpdate(roomId string, roomState entities.Room) {
	h.Hub.BroadcastToRoom(BroadcastMessage{
		RoomID: roomId,
		Message: Message{
			Type:    "ROOM_STATE_UPDATED",
			Payload: utils.MustMarshal(roomState),
		},
	})
}