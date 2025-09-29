package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"planning-poker/internal/infrastructure/persistence"
	"planning-poker/internal/ports/repository"

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

type Client struct {
	Conn   *websocket.Conn
	Hub    Hub 
	RoomID string
	UserID string
	Ctx    context.Context
	Redis repository.RedisRepository
}

type BroadcastMessage struct {
	RoomID  string
	Message Message
}

type Hub interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastToRoom(message BroadcastMessage)
}


type HttpWsHandler struct {
	Hub             Hub
	RedisRepository repository.RedisRepository
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewWsHandler(hub Hub, redisClient *redis.Client) *HttpWsHandler {
	return &HttpWsHandler{
		Hub:             hub,
		RedisRepository: persistence.NewRedisRepositoryImpl(redisClient),
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
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		Conn:   conn,
		Hub:    h.Hub,
		RoomID: roomId,
		UserID: userId,
		Ctx:   ctx,
		Redis: h.RedisRepository,
	}

	// Registrar o client no hub
	h.Hub.RegisterClient(client)

	// Iniciar a leitura de mensagens do WebSocket
	go func() {
		defer func() {
			cancel()
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

			log.Printf("Received message: %+v", msg)
		}
	}()
}