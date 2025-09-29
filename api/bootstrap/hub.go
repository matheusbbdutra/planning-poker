package bootstrap

import (
	"log"
	"planning-poker/internal/http/handler"
)

// Hub mantém o conjunto de clientes ativos e transmite mensagens para eles.
type Hub struct {
	// Clientes registrados. A chave é o ID da sala (roomId).
	rooms map[string]map[*handler.Client]bool

	// Mensagens de entrada dos clientes.
	broadcast chan handler.BroadcastMessage

	// Registrar requisições do hub.
	register chan *handler.Client

	// Desregistrar requisições.
	unregister chan *handler.Client
}

func NewHub() *Hub {
	hub := &Hub{
		register:   make(chan *handler.Client),
		broadcast:  make(chan handler.BroadcastMessage),
		unregister: make(chan *handler.Client),
		rooms:      make(map[string]map[*handler.Client]bool),
	}
	go hub.run()
	log.Println("WebSocket Hub criado e rodando.")
	return hub
}

// RegisterClient implementa o primeiro método da nossa interface.
// Ele simplesmente envia o cliente para o canal de registro.
func (h *Hub) RegisterClient(client *handler.Client) {
	h.register <- client
}

// UnregisterClient implementa o segundo método da nossa interface.
func (h *Hub) UnregisterClient(client *handler.Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastToRoom(message handler.BroadcastMessage) {
	h.broadcast	<- message
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*handler.Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			log.Printf("Cliente registrado na sala %s. Total de clientes na sala: %d", client.RoomID, len(h.rooms[client.RoomID]))
			roomState, err := client.Redis.GetRoom(client.Ctx, client.RoomID)

			if err != nil {
				log.Printf("Erro ao buscar estado inicial da sala %s: %v", client.RoomID, err)
				break
			}
			if roomState == nil {
				log.Printf("Sala %s não encontrada no Redis.", client.RoomID)
				break
			}
			log.Printf("Estado inicial da sala %s: %+v", client.RoomID, roomState)

			initialState := handler.Message{
				Type:    "INITIAL_STATE",
				Payload: roomState,
			}
			if err := client.Conn.WriteJSON(initialState); err != nil {
				log.Printf("Erro ao enviar estado inicial para o cliente %s: %v", client.UserID, err)
				client.Conn.Close()
			}

		case client := <-h.unregister:
			if _, ok := h.rooms[client.RoomID]; ok {
				delete(h.rooms[client.RoomID], client)
				client.Conn.Close()
				if len(h.rooms[client.RoomID]) == 0 {
					delete(h.rooms, client.RoomID)
				}
				log.Printf("Cliente %s removido da sala %s", client.UserID, client.RoomID)
			}

		case message := <-h.broadcast:
			if clients, ok := h.rooms[message.RoomID]; ok {
				for client := range clients {
					err := client.Conn.WriteJSON(message.Message)
					if err != nil {
						log.Printf("Erro ao enviar mensagem para o cliente %s: %v", client.UserID, err)
						client.Conn.Close()
						delete(clients, client)
					}
				}
			}
		}
	}
}
