package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"GameTicTacToe/dbs"
	"GameTicTacToe/game"
)

type result struct {
	Result string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	conn := dbs.DbConn{}
	conn.InitDb()
	id, er := r.Cookie("UserId")
	if er != nil {
		ck := &http.Cookie{
			Name:    "UserId",
			Value:   game.GenUsedCk(),
			Expires: time.Now(),
			MaxAge:  9000,
		}
		http.SetCookie(w, ck)
		t, _ := template.ParseFiles("pages/registration.html")
		t.Execute(w, nil)
		_, err := conn.CreateUser(ck.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		name := r.FormValue("name")
		if name != "" {
			err := conn.AddName(id.Value, name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if name == "" {
			name, _ = conn.GetUserName(id.Value)
		}
		res := result{Result: name}
		t, _ := template.ParseFiles("pages/home.html")
		t.Execute(w, res)
	}
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	conn := dbs.DbConn{}
	conn.InitDb()
	gameData := &dbs.GameData{}
	gi, err := r.Cookie("GameId")
	if err != nil {
		id, _ := r.Cookie("UserId")
		gameData, err = conn.CreateNewGame(id.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ck := &http.Cookie{
			Name:  "GameId",
			Value: strconv.Itoa(gameData.GameId),
		}
		http.SetCookie(w, ck)
		//wait to second player
	} else {
		gameData, _ = conn.RefreshGameData(gi.Value)
	}
	gameData.Who = "You can make move!"
	t, _ := template.ParseFiles("pages/game.html")
	t.Execute(w, gameData)
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	conn := dbs.DbConn{}
	conn.InitDb()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	gameId := r.FormValue("gameId")
	id, _ := r.Cookie("UserId")
	players, _ := conn.GetPlayersCK(gameId)
	if players.PlayerOId == 0 {
		err := conn.SetPlayerId(id.Value, gameId)
		idO, _ := strconv.Atoi(id.Value)
		players.PlayerOId = idO
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if conn.CheckGame(gameId) && ((strconv.Itoa(players.PlayerOId) == id.Value) || (strconv.Itoa(players.PlayerXId) == id.Value)) {
		ck := &http.Cookie{
			Name:  "GameId",
			Value: gameId,
		}
		http.SetCookie(w, ck)
		http.Redirect(w, r, "/game", 301)

		conn.SetGameStatus(gameId, dbs.Running)
		gameData, _ := conn.RefreshGameData(gameId)
		t, _ := template.ParseFiles("pages/game.html")
		t.Execute(w, gameData)
	} else {
		http.Redirect(w, r, "/", 301)
		name, _ := conn.GetUserName(id.Value)
		res := result{Result: name}
		t, _ := template.ParseFiles("pages/home.html")
		t.Execute(w, res)
	}
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
		http.Redirect(w, r, "/", 301)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, _, err := conn.GetGameData(gameId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	players, err := conn.GetPlayersCK(gameId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := r.Cookie("UserId")
	ck, _ := strconv.Atoi(id.Value)
	if (data.Symbol == "X" && (players.PlayerXId == ck)) || (data.Symbol == "O" && (players.PlayerOId == ck)) {
		gameData, _, count, err := conn.MakeMove(numOfCell, gameId.Value) // you can use for unique
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
			var name string
			if win == "X" {
				name, _ = conn.GetUserName(strconv.Itoa(players.PlayerXId))
			} else {
				name, _ = conn.GetUserName(strconv.Itoa(players.PlayerOId))
			}
			res := result{
				Result: fmt.Sprintf("Winner is %s", name),
			}
			ck := &http.Cookie{
				Name:   "GameId",
				MaxAge: -1,
			}
			http.SetCookie(w, ck)
			t, _ := template.ParseFiles("pages/finish.html")
			t.Execute(w, res)
			conn.SetGameStatus(gameId.Value, dbs.Finished)
		} else if count >= 10 {
			res := result{
				Result: "Nobody win!",
			}
			ck := &http.Cookie{
				Name:   "GameId",
				MaxAge: -1,
			}
			http.SetCookie(w, ck)
			t, _ := template.ParseFiles("pages/finish.html")
			t.Execute(w, res)
			conn.SetGameStatus(gameId.Value, dbs.Finished)
		} else {
			// checking if a move has been made?
			if data.Symbol != gameData.Symbol {
				if gameData.Symbol == "O" {
					name, _ := conn.GetUserName(strconv.Itoa(players.PlayerOId))
					gameData.Who = fmt.Sprintf("%s is making move!", name)
					gameData.Symbol = "X"
				} else {
					name, _ := conn.GetUserName(strconv.Itoa(players.PlayerXId))
					gameData.Who = fmt.Sprintf("%s making move!", name)
					gameData.Symbol = "O"
				}
			} else {
				gameData.Who = "You can make move!"
			}
			t, _ := template.ParseFiles("pages/game.html")
			t.Execute(w, gameData)
		}
	} else {
		if data.Symbol == "O" {
			name, _ := conn.GetUserName(strconv.Itoa(players.PlayerOId))
			data.Who = fmt.Sprintf("%s is making move!", name)
			data.Symbol = "X"
		} else {
			name, _ := conn.GetUserName(strconv.Itoa(players.PlayerXId))
			data.Who = fmt.Sprintf("%s is making move!", name)
			data.Symbol = "O"
		}
		t, _ := template.ParseFiles("pages/game.html")
		t.Execute(w, data)
	}
}
