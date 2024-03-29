package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	database "rest-websockets/database"
	repository "rest-websockets/repository"
	websocket "rest-websockets/websocket"
)

// Config Server
type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

// Interface Server
type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

// Server
type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

// Implements config methods of interface to return config server
func (b *Broker) Config() *Config {
	return b.config
}

// Create new server
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

// Create new router
func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	handler := cors.Default().Handler(b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}
