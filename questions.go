package main

import (
	"encoding/json"
	"log"
	"os"
)

func loadQuestions() Quiz {
	var quiz Quiz
	file, err := os.ReadFile("questions.json")
	if err != nil {
		log.Fatal("Unable to load question file: ", err)
	}
	err = json.Unmarshal(file, &quiz)
	if err != nil {
		log.Fatal("Unable to unmarshal question file: ", err)
	}
	
	log.Println("Quiz loaded successfully.")
	return quiz
}

