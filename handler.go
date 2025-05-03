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
      case "restart_game":
      	handleRestartGame(conn, msg)
      case "list_rooms":
      	handleListRooms(conn)
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
		GameOver: false,
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

// func handleNextQuestion(conn *websocket.Conn, msg Message) {
// 	quizHub.Mu.Lock()
// 	defer quizHub.Mu.Unlock()
// 
// 	room, ok := quizHub.Rooms[msg.RoomID]
// 	if !ok {
// 		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
// 		return
// 	}
// 
// 	log.Println("CurrentQ: %d >= Question: %d", room.CurrentQ, len(room.Questions))
// 	if room.CurrentQ >= len(room.Questions) {
// 		// End game automatically if no more questions
// 		handleEndGame(conn, msg)
// 		return
// 	}
// 
// 	question := room.Questions[room.CurrentQ]
// 	payload := map[string]interface{}{
// 		"type":     "question",
// 		"tag":      question.Tag,
// 		"question": question.Question,
// 		"options": map[string]string{
// 			"a": question.OptionA,
// 			"b": question.OptionB,
// 			"c": question.OptionC,
// 			"d": question.OptionD,
// 		},
// 	}
// 
// 	for player := range room.Players {
// 		player.WriteJSON(payload)
// 	}
// 
// 	room.CurrentQ++
// }

func handleNextQuestion(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	// Check if the room has been closed
  if room.GameOver {
		conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "Game is over. Use 'restart_game' to play again.",
		})
		return
	}

	// End game if there are no more questions
	if room.CurrentQ >= len(room.Questions) {
		room.GameOver = true
		// Send "game_over" to all players only once
		payload := map[string]string{
			"type":    "game_over",
			"message": "No more questions. Game over!",
		}
		for player := range room.Players {
			player.WriteJSON(payload)
		}
		return
	}

	// Send the next question
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

	room.CurrentQ++ // Increment after sending
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

func handleRestartGame(conn *websocket.Conn, msg Message) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	room, ok := quizHub.Rooms[msg.RoomID]
	if !ok {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	room.CurrentQ = 0
	room.GameOver = false

	if len(room.Questions) == 0 {
		conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "No questions available.",
		})
		return
	}

	first := room.Questions[0]
	payload := map[string]interface{}{
		"type":     "question",
		"tag":      first.Tag,
		"question": first.Question,
		"options": map[string]string{
			"a": first.OptionA,
			"b": first.OptionB,
			"c": first.OptionC,
			"d": first.OptionD,
		},
	}

	for player := range room.Players {
		player.WriteJSON(payload)
	}

	room.CurrentQ++ // Prep for next
}


func handleListRooms(conn *websocket.Conn) {
	quizHub.Mu.Lock()
	defer quizHub.Mu.Unlock()

	var roomList []map[string]interface{}
	for id, room := range quizHub.Rooms {
  	roomList = append(roomList, map[string]interface{}{
  		"room_id":          id,
  		"player_count":     len(room.Players),
  		"game_over":        room.GameOver,
  		"current_question": room.CurrentQ,
  		"total_questions":  len(room.Questions),
  	})
  }

	response := map[string]interface{}{
		"type":  "room_list",
		"rooms": roomList,
	}

	if err := conn.WriteJSON(response); err != nil {
	  log.Printf("error writing list_rooms response: %v", err)
  }

}
