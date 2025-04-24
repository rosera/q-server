package main

import (
//	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read error:", err)
			break
		}

		switch msg.Type {
      case "create_room":
      	handleCreateRoom(conn, msg)
      case "join_room":
      	handleJoinRoom(conn, msg)
      case "start_game":
      	handleStartGame(conn, msg)
      case "next_question":
      	handleNextQuestion(conn, msg)
      case "end_game":
      	handleEndGame(conn, msg)
    }

	}
}

func handleCreateRoom(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	if _, exists := quizHub.Rooms[msg.RoomID]; exists {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room already exists"})
		return
	}

	quiz := loadQuestions() // Load from file
	quizHub.Rooms[msg.RoomID] = &Room{
		ID:       msg.RoomID,
		Admin:    conn,
		Players:  make(map[*websocket.Conn]string),
		Questions: quiz.Tasks,
		CurrentQ: 0,
	}
	conn.WriteJSON(map[string]string{"type": "room_created", "room_id": msg.RoomID})
}

func handleJoinRoom(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}
	room.Players[conn] = msg.Name
	conn.WriteJSON(map[string]string{"type": "joined", "room_id": msg.RoomID})
}

func handleStartGame(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}
	question := room.Questions[room.CurrentQ]
	payload := map[string]interface{}{
		"type":     "question",
		"tag":      question.Tag,
		"question": question.Question,
		"options": map[string]string{
			"a": question.OptionA,
			"b": question.OptionB,
			"c": question.OptionC,
			"d": question.OptionD,
		},
	}
	for player := range room.Players {
		player.WriteJSON(payload)
	}
}

func handleNextQuestion(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	room.CurrentQ++
	if room.CurrentQ >= len(room.Questions) {
		// End game automatically if no more questions
		handleEndGame(conn, msg)
		return
	}

	question := room.Questions[room.CurrentQ]
	payload := map[string]interface{}{
		"type":     "question",
		"tag":      question.Tag,
		"question": question.Question,
		"options": map[string]string{
			"a": question.OptionA,
			"b": question.OptionB,
			"c": question.OptionC,
			"d": question.OptionD,
		},
	}
	for player := range room.Players {
		player.WriteJSON(payload)
	}
}

func handleEndGame(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	// Reset the question index to restart quiz from the beginning
	room.CurrentQ = 0

	payload := map[string]string{
		"type":    "end_game",
		"message": "Game Over! Thanks for playing.",
	}
	for player := range room.Players {
		player.WriteJSON(payload)
	}
}

