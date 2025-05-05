package main

import (
	"log"
	"fmt"
	"net/http"
)

var (
  hostIp = "192.168.87.55"
  hostPort = "8080"
)

func main() {

	hostInformation := fmt.Sprintf("%s:%s", hostIp, hostPort)

	http.HandleFunc("/ws", handleConnections)
	log.Println("Quiz server started on:", hostInformation)
	err := http.ListenAndServe(hostInformation, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
