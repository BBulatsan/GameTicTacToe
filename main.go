package main

import (
	"fmt"
	"log"
	"net/http"

	"GameTicTacToe/handlers"
)

const (
	port = ":5000"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/new_game", handlers.StartHandler)
	http.HandleFunc("/connect", handlers.ConnectHandler)
	http.HandleFunc("/game", handlers.GameHandler)
	fmt.Println("Server is running on port" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
