package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Player struct {
	Name string `json:"name"`
}

type Room struct {
	ID       string
	Admin    *websocket.Conn
	Players  map[*websocket.Conn]string
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
