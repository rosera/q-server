// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Question struct {
	Tag      string `json:"tag"`
	Question string `json:"question"`
	OptionA  string `json:"option_a"`
	OptionB  string `json:"option_b"`
	OptionC  string `json:"option_c"`
	OptionD  string `json:"option_d"`
	Answer   string `json:"answer"`
}

type Quiz struct {
	ID     string     `json:"id"`
	Tasks  []Question `json:"tasks"`
	Author string     `json:"author"`
}

type Player struct {
	Name       string
	Conn       *websocket.Conn
	Score      int
	AnswerChan chan string
}

type Room struct {
	sync.Mutex
	Players      []*Player
	Quiz         Quiz
	CurrentIndex int
	Answers      map[string]string
	StartRequested bool
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var rooms = make(map[string]*Room)
var roomLock = sync.Mutex{}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	var initMsg struct {
		PlayerName string `json:"playerName"`
		RoomID     string `json:"roomID"`
		Role       string `json:"role"` // "admin" or "player"
	}

	if err := conn.ReadJSON(&initMsg); err != nil {
		log.Println("Init message error:", err)
		return
	}

	room := getOrCreateRoom(initMsg.RoomID)

	// isAdmin := initMsg.Role == "admin"
	player := &Player{
		Name:       initMsg.PlayerName,
		Conn:       conn,
		AnswerChan: make(chan string),
	}
	room.Lock()
	room.Players = append(room.Players, player)
	room.Unlock()

	broadcastToRoom(room, map[string]string{
		"type": "player_joined",
		"name": player.Name,
	})

	go handlePlayer(room, player)
	if len(room.Players) == 2 {
		startGame(room)
	}
}

func getOrCreateRoom(id string) *Room {
	roomLock.Lock()
	defer roomLock.Unlock()
	if room, exists := rooms[id]; exists {
		return room
	}
	quiz := loadSampleQuiz()
	room := &Room{
		Quiz:    quiz,
		Answers: make(map[string]string),
	}
	rooms[id] = room
	return room
}

func handlePlayer(room *Room, player *Player) {
	defer func() {
		room.Lock()
		// Remove player
		for i, p := range room.Players {
			if p == player {
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
				break
			}
		}
		delete(room.Answers, player.Name)
		room.Unlock()
		broadcastToRoom(room, map[string]string{
			"type": "player_left",
			"name": player.Name,
		})
		player.Conn.Close()
	}()

	for {
  	var msg map[string]interface{}
  	if err := player.Conn.ReadJSON(&msg); err != nil {
  		log.Printf("Player %s disconnected: %v\n", player.Name, err)
  		return
  	}
  
  	msgType := msg["type"]
  	switch msgType {
  	case "answer":
  		playerName := msg["playerName"].(string)
  		answer := msg["answer"].(string)
  		room.Lock()
  		room.Answers[playerName] = answer
  		room.Unlock()
  
  	case "start_game":
  		room.Lock()
  		if !room.StartRequested {
  			room.StartRequested = true
  			go startGame(room)
  		}
  		room.Unlock()
  	}
  }

//	for {
//		var msg struct {
//			PlayerName string `json:"playerName"`
//			Answer     string `json:"answer"`
//		}
//		err := player.Conn.ReadJSON(&msg)
//		if err != nil {
//			log.Printf("Player %s disconnected: %v\n", player.Name, err)
//			return
//		}
//		room.Lock()
//		// room.Answers[player.Name] = msg.Answer
//		room.Answers[msg.PlayerName] = msg.Answer
//		room.Unlock()
//	}
}

func startGame(room *Room) {
	sendQuestionToAll(room)

	go func() {
		for {
			room.Lock()
			if len(room.Answers) < len(room.Players) {
				room.Unlock()
				continue
			}

			current := room.Quiz.Tasks[room.CurrentIndex]
			for _, p := range room.Players {
				if room.Answers[p.Name] == current.Answer {
					p.Score++
				}
			}

			room.Answers = make(map[string]string)
			room.CurrentIndex++
			if room.CurrentIndex >= len(room.Quiz.Tasks) {
				declareWinner(room)
				room.Unlock()
				return
			}
			sendQuestionToAll(room)
			room.Unlock()
		}
	}()
}

func sendQuestionToAll(room *Room) {
	current := room.Quiz.Tasks[room.CurrentIndex]
	msg := map[string]interface{}{
		"type":     "question",
		"question": current,
	}
	broadcastToRoom(room, msg)
}

func declareWinner(room *Room) {
	var winner *Player
	highScore := -1
	for _, p := range room.Players {
		if p.Score > highScore {
			winner = p
			highScore = p.Score
		}
	}
	broadcastToRoom(room, map[string]interface{}{
		"type":   "end",
		"winner": winner.Name,
		"score":  winner.Score,
	})
}

func broadcastToRoom(room *Room, msg interface{}) {
	for _, p := range room.Players {
		_ = p.Conn.WriteJSON(msg)
	}
}

func loadSampleQuiz() Quiz {
	raw := `
	{
	  "author": "cmdlinezero",
	  "id": "General Knowledge",
	  "tasks": [
		{
		  "tag": "1",
		  "question": "What is the capital city of Australia?",
		  "option_a": "Sydney",
		  "option_b": "Melbourne",
		  "option_c": "Canberra",
		  "option_d": "Perth",
		  "answer": "option_c"
		}
	  ]
	}`
	var q Quiz
	_ = json.Unmarshal([]byte(raw), &q)
	return q
}

