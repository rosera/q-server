package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Player struct {
	Name string `json:"name"`
	ClientID string `json:"client_id"`
}

type Room struct {
	ID       string
	Admin    *websocket.Conn
	// Players  map[*websocket.Conn]string
	Players  map[*websocket.Conn]Player
	PlayerList []Player
	Questions []Task
	CurrentQ int
	GameOver bool
}

type Hub struct {
	Rooms map[string]*Room
	Mu    sync.Mutex
}

var quizHub = Hub{
	Rooms: make(map[string]*Room),
}
