package bootstrap

import "log"

// Client representa uma única conexão de websocket.
type Client struct {
	// Adicionaremos mais campos aqui, como a conexão em si.
}

// Hub mantém o conjunto de clientes ativos e transmite mensagens para eles.
type Hub struct {
	// Clientes registrados. A chave é o ID da sala (roomId).
	rooms map[string]map[*Client]bool

	// Mensagens de entrada dos clientes.
	Broadcast chan []byte

	// Registrar requisições do hub.
	register chan *Client

	// Desregistrar requisições.
	unregister chan *Client
}

func NewHub() *Hub {
	// No futuro, aqui teremos a lógica para iniciar o Hub em uma goroutine.
	log.Println("WebSocket Hub criado.")
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}
