package bootstrap

import (
	"log"
	"planning-poker/internal/http/handler"
)

type App struct {
	Server      *Server
	RoomHandler *handler.RoomHandler
	WsHandler   *handler.HttpWsHandler
	WsHub       *Hub
}

func NewApp() *App {
	a := &App{
		Server: NewServer(),
		WsHub:  NewHub(),
	}
	a.initHandlers()
	return a
}

func (a *App) Start() {

	if err := a.Server.InitRedis(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	a.initHandlers()
	a.Server.Start(a.RoomHandler, a.WsHandler)
}

func (a *App) initHandlers() {
	a.RoomHandler = handler.NewRoomHandler(a.Server.redis)
	a.WsHandler = handler.NewWsHandler(a.WsHub, a.Server.redis)
}
