package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Room struct {
	ID       string
	Admin    *websocket.Conn
	Players  map[*websocket.Conn]string
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

