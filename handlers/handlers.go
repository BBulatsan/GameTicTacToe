package handlers

import (
	"encoding/json"
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

type Handler struct {
	Conn dbs.DbConn
}

const (
	userID = "user_id"
	gameID = "game_id"
)

// HomeHandler use for present first page and registration.
func (h Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	ck := &http.Cookie{
		Name:   gameID,
		MaxAge: -1,
	}
	http.SetCookie(w, ck)
	id, err := r.Cookie(userID)
	if err != nil {
		ck = &http.Cookie{
			Name:   userID,
			Value:  game.GenUsedCk(),
			MaxAge: 2147483647,
		}
		http.SetCookie(w, ck)
		t, _ := template.ParseFiles("pages/registration.html")
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = h.Conn.CreateUser(ck.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		if err = r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		name := r.FormValue("name")
		if name != "" {
			err = h.Conn.AddName(id.Value, name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if name == "" {
			name, _ = h.Conn.GetUserName(id.Value)
		}
		res := result{Result: name}
		t, _ := template.ParseFiles("pages/home.html")
		err = t.Execute(w, res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// StartHandler use for create new game.
func (h Handler) StartHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(gameID)
	if err != nil {
		gameData := &dbs.GameData{
			Symbol: game.GenSymbol(),
			One:    "1",
			Two:    "2",
			Three:  "3",
			Four:   "4",
			Five:   "5",
			Six:    "6",
			Seven:  "7",
			Eight:  "8",
			Nine:   "9",
		}

		gId, err := h.Conn.CreateGame()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, _ := r.Cookie(userID)
		gameId := strconv.Itoa(gId)

		err = h.Conn.SetPlayerId(id.Value, gameId, gameData.Symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gameData.GameId = gId
		gd, err := json.Marshal(gameData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.Conn.CreateMove(gd, 1, gameId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ck := &http.Cookie{
			Name:  gameID,
			Value: gameId,
		}
		http.SetCookie(w, ck)
		http.Redirect(w, r, "/game", http.StatusMovedPermanently)
		//TODO wait to second player
		return
	} else {
		http.Redirect(w, r, "/connect", http.StatusMovedPermanently)
		return
	}
}

// ConnectHandler use for connect to new games.
func (h Handler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	gameId := r.FormValue(gameID)
	status, err := h.Conn.GetGameStatus(gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if status == dbs.Finished {
		res := result{
			Result: "This game has been finished!",
		}
		t, _ := template.ParseFiles("pages/finish.html")
		err = t.Execute(w, res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// This part check players in game and create a pair of players.
	id, _ := r.Cookie(userID)
	players, _ := h.Conn.GetPlayersCK(gameId)
	idValue, _ := strconv.Atoi(id.Value)
	if players.PlayerOId == 0 && players.PlayerXId != idValue {
		err = h.Conn.SetPlayerId(id.Value, gameId, dbs.O)
		players.PlayerOId = idValue
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if players.PlayerXId == 0 && players.PlayerOId != idValue {
		err = h.Conn.SetPlayerId(id.Value, gameId, dbs.X)
		players.PlayerXId = idValue
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if h.Conn.CheckGame(gameId) && ((strconv.Itoa(players.PlayerOId) == id.Value) || (strconv.Itoa(players.PlayerXId) == id.Value)) {
		ck := &http.Cookie{
			Name:  gameID,
			Value: gameId,
		}
		http.SetCookie(w, ck)
		http.Redirect(w, r, "/game", http.StatusMovedPermanently)

		err = h.Conn.SetGameStatus(gameId, dbs.Running)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gameData, _ := h.Conn.RefreshGameData(gameId)
		t, _ := template.ParseFiles("pages/game.html")
		err = t.Execute(w, gameData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Redirect(w, r, "/", 301)
		name, _ := h.Conn.GetUserName(id.Value)
		res := result{Result: name}
		t, _ := template.ParseFiles("pages/home.html")
		err = t.Execute(w, res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//GameHandler use for present game and implement main logic.
func (h Handler) GameHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	numOfCell := r.FormValue("numOfCell")

	gameId, err := r.Cookie(gameID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data, _, err := h.Conn.GetGameData(gameId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	players, err := h.Conn.GetPlayersCK(gameId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := r.Cookie(userID)
	ck, _ := strconv.Atoi(id.Value)

	// checking who make move
	if (data.Symbol == "X" && (players.PlayerXId == ck)) || (data.Symbol == "O" && (players.PlayerOId == ck)) {
		g := game.NewGame(h.Conn)
		gameData, _, count, err := g.MakeMove(numOfCell, gameId.Value) // you can use for unique
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
				name, _ = h.Conn.GetUserName(strconv.Itoa(players.PlayerXId))
			} else {
				name, _ = h.Conn.GetUserName(strconv.Itoa(players.PlayerOId))
			}
			res := result{
				Result: fmt.Sprintf("Winner is %s", name),
			}
			ck := &http.Cookie{
				Name:   gameID,
				MaxAge: -1,
			}
			http.SetCookie(w, ck)
			t, _ := template.ParseFiles("pages/finish.html")
			err = t.Execute(w, res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = h.Conn.SetGameStatus(gameId.Value, dbs.Finished)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if count >= 10 {
			res := result{
				Result: "Nobody win!",
			}
			ck := &http.Cookie{
				Name:   gameID,
				MaxAge: -1,
			}
			http.SetCookie(w, ck)
			t, _ := template.ParseFiles("pages/finish.html")
			err = t.Execute(w, res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = h.Conn.SetGameStatus(gameId.Value, dbs.Finished)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// checking if a move has been made?
			if data.Symbol != gameData.Symbol {
				if gameData.Symbol == "O" {
					name, _ := h.Conn.GetUserName(strconv.Itoa(players.PlayerOId))
					gameData.Who = fmt.Sprintf("%s is making move!", name)
					gameData.Symbol = "X"
				} else {
					name, _ := h.Conn.GetUserName(strconv.Itoa(players.PlayerXId))
					gameData.Who = fmt.Sprintf("%s making move!", name)
					gameData.Symbol = "O"
				}
			} else {
				gameData.Who = "You can make move!"
			}
			t, _ := template.ParseFiles("pages/game.html")
			err = t.Execute(w, gameData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		if data.Symbol == "O" {
			name, _ := h.Conn.GetUserName(strconv.Itoa(players.PlayerOId))
			data.Who = fmt.Sprintf("%s is making move!", name)
			data.Symbol = "X"
		} else {
			name, _ := h.Conn.GetUserName(strconv.Itoa(players.PlayerXId))
			data.Who = fmt.Sprintf("%s is making move!", name)
			data.Symbol = "O"
		}
		t, _ := template.ParseFiles("pages/game.html")
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
