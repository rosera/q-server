package main

type Quiz struct {
	Author   string `json:"author"`
	Endpoint string `json:"endpoint"`
	ID       string `json:"id"`
	Publish  string `json:"publish"`
	Tasks    []Task `json:"tasks"`
	URI      string `json:"uri"`
}

type Task struct {
	Tag      string `json:"tag"`
	Question string `json:"question"`
	OptionA  string `json:"option_a"`
	OptionB  string `json:"option_b"`
	OptionC  string `json:"option_c"`
	OptionD  string `json:"option_d"`
	Answer   string `json:"answer"`
}

type Message struct {
	Type   string `json:"type"`
	RoomID string `json:"room_id"`
	Name   string `json:"name,omitempty"`
}

