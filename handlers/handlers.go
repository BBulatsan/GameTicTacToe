package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"GameTicTacToe/dbs"
	"GameTicTacToe/game"
)

type result struct {
	Result string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pages/home.html")
	t.Execute(w, nil)
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	conn := dbs.DbConn{}
	conn.InitDb()
	gameData, err := conn.CreateNewGame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ck := &http.Cookie{
		Name:  "GameId",
		Value: strconv.Itoa(gameData.GameId),
	}
	http.SetCookie(w, ck)

	t, _ := template.ParseFiles("pages/game.html")
	t.Execute(w, gameData)
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	gameId := r.FormValue("gameId")
	ck := &http.Cookie{
		Name:  "GameId",
		Value: gameId,
	}
	http.SetCookie(w, ck)

	fmt.Fprintf(w, "Id = %s\n", gameId)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	conn := dbs.DbConn{}
	conn.InitDb()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	numOfCell := r.FormValue("numOfCell")

	gameId, err := r.Cookie("GameId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	gameData, count, err := conn.MakeMove(numOfCell, gameId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	finish, win, err := game.CheckWin(gameData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if finish {
		res := result{
			Result: fmt.Sprintf("Winer is %s", win),
		}

		t, _ := template.ParseFiles("pages/finish.html")
		t.Execute(w, res)
		conn.SetFinish(gameId.Value)
	} else if count >= 10 {
		res := result{
			Result: "Nobody win!",
		}
		t, _ := template.ParseFiles("pages/finish.html")
		t.Execute(w, res)
		conn.SetFinish(gameId.Value)
	} else {
		t, _ := template.ParseFiles("pages/game.html")
		t.Execute(w, gameData)
	}
}
