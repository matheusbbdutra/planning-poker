package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"planning-poker/internal/http/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	client *chi.Mux
	redis  *redis.Client
}

func NewServer() *Server {
	return &Server{
		client: chi.NewRouter(),
	}
}

func (s *Server) Start(roomHandler *handler.RoomHandler, wsHandler *handler.HttpWsHandler) {
	r := initRoutes(roomHandler, wsHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("WebSocket server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server exiting")
}

func initRoutes(roomHandler *handler.RoomHandler, wsHandler *handler.HttpWsHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/room", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				if err := roomHandler.NewRoom(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})
			r.Post("/{roomId}/join", func(w http.ResponseWriter, r *http.Request) {
				if err := roomHandler.JoinRoom(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})
		})
		r.Get("/room/ws/{roomId}/{userId}", wsHandler.HandleWebSocket)
	})
	return r
}

func (s *Server) InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	log.Printf("Connected to Redis successfully")

	s.redis = redisClient
	return nil
}
