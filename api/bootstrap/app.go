package bootstrap

import (
	"log"
	httphandler "planning-poker/internal/http"
)

type App struct {
	Server      *Server
	RoomHandler *httphandler.HttpRoomHandler
	WsHub       *Hub
}

func NewApp() *App {
	a := &App{Server: NewServer()}
	a.initHandlers()
	return a
}

func (a *App) Start() {

	if err := a.Server.InitRedis(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	a.initHandlers()
	a.Server.Start(a.RoomHandler)
}

func (a *App) initHandlers() {
	a.RoomHandler = httphandler.NewRoomHandler(a.Server.redis)
}
