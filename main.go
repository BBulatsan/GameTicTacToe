package main

import (
	"GameTicTacToe/dbs"
	"fmt"
	"log"
	"net/http"

	"GameTicTacToe/handlers"
)

const (
	port = ":5000"
)

func main() {
	conn := dbs.DbConn{}
	err := conn.InitDb()
	if err != nil {
		log.Fatal("Error connect to DB", err)
		return
	}
	defer conn.CloseDb()

	handler := handlers.Handler{Conn: conn}
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/new_game", handler.StartHandler)
	http.HandleFunc("/connect", handler.ConnectHandler)
	http.HandleFunc("/game", handler.GameHandler)
	fmt.Println("Server is running on port" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
