package dbs

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//CREATE TABLE IF NOT EXISTS games (id INTEGER PRIMARY KEY, status message_text, count_players SMALLINT, player_x_id INTEGER, player_o_id INTEGER);
//CREATE TABLE IF NOT EXISTS moves (game_id INTEGER, map_of_moves json, who_move message_text, FOREIGN KEY (game_id)  REFERENCES games (id) ON DELETE CASCADE);
//CREATE TABLE IF NOT EXISTS users ( id  integer constraint user_pk primary key, name text );
const (
	new      = "new"
	running  = "running"
	finished = "finished"
)

type games struct {
	Id           int
	Status       string
	CountPlayers int
}

type moves struct {
	GameId    int
	GameData  string
	CountMove int
}

type GameData struct {
	Who    string
	GameId int
	One    string
	Two    string
	Three  string
	Four   string
	Five   string
	Six    string
	Seven  string
	Eight  string
	Nine   string
}

type DbConn struct {
	conn *sql.DB
}

func (d *DbConn) InitDb() {
	db, err := sql.Open("sqlite3", "./dbs/store/store.db")
	if err != nil {
		panic(err)
	}
	d.conn = db
}

func (d *DbConn) CreateNewGame() (*GameData, error) {
	var gameId int
	moveMap := &GameData{
		Who:   "X",
		One:   "1",
		Two:   "2",
		Three: "3",
		Four:  "4",
		Five:  "5",
		Six:   "6",
		Seven: "7",
		Eight: "8",
		Nine:  "9",
	}
	statement := "INSERT INTO games (status, count_players) VALUES ($1, $2);"
	res, err := d.conn.Exec(statement, new, 1)
	if err != nil {
		return moveMap, err
	}
	id, err := res.LastInsertId()
	gameId = int(id)
	moveMap.GameId = gameId

	js, err := json.Marshal(moveMap)
	if err != nil {
		return moveMap, err
	}
	statement = "INSERT INTO moves (game_id, game_data, count_move) VALUES ($1, $2, $3);"
	_, err = d.conn.Exec(statement, gameId, js, 1)
	if err != nil {
		return moveMap, err
	}

	return moveMap, nil
}

func (d *DbConn) MakeMove(move string, gameId string) (*GameData, int, error) {
	gameData := &GameData{}
	GData, count, err := d.getGameData(gameId)
	if err != nil {
		return gameData, count, err
	}
	err = json.Unmarshal([]byte(GData), gameData)
	if err != nil {
		return gameData, count, err
	}
	if gameData.Who == "X" {
		gameData, err = executeGame(gameData, move, "X")
		if err != nil {
			return gameData, count, nil
		}
		gameData.Who = "O"
	} else {
		gameData, err = executeGame(gameData, move, "O")
		if err != nil {
			return gameData, count, nil
		}
		gameData.Who = "X"
	}

	gd, err := json.Marshal(gameData)
	if err != nil {
		return gameData, count, err
	}

	statement := "UPDATE moves SET game_data=$1, count_move=$2 WHERE game_id=$3;"
	count++
	_, err = d.conn.Exec(statement, gd, count, gameId)
	if err != nil {
		count--
		return gameData, count, err
	}

	return gameData, count, nil
}

func (d *DbConn) getGameData(gameId string) (string, int, error) {
	move := &moves{}
	rows, err := d.conn.Query(fmt.Sprintf("SELECT game_data, count_move FROM moves WHERE game_id=%s", gameId))
	if err != nil {
		return move.GameData, move.CountMove, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&move.GameData, &move.CountMove)
	if err != nil {
		return move.GameData, move.CountMove, err
	}

	return move.GameData, move.CountMove, nil
}

func (d *DbConn) SetFinish(gameId string) error {
	statement := "UPDATE games SET status=$1 WHERE id=$2;"
	_, err := d.conn.Exec(statement, finished, gameId)
	if err != nil {
		return err
	}
	return nil
}

func executeGame(gameData *GameData, move string, symbol string) (*GameData, error) {
	if gameData.One == move {
		gameData.One = symbol
		return gameData, nil
	}
	if gameData.Two == move {
		gameData.Two = symbol
		return gameData, nil
	}
	if gameData.Three == move {
		gameData.Three = symbol
		return gameData, nil
	}
	if gameData.Four == move {
		gameData.Four = symbol
		return gameData, nil
	}
	if gameData.Five == move {
		gameData.Five = symbol
		return gameData, nil
	}
	if gameData.Six == move {
		gameData.Six = symbol
		return gameData, nil
	}
	if gameData.Seven == move {
		gameData.Seven = symbol
		return gameData, nil
	}
	if gameData.Eight == move {
		gameData.Eight = symbol
		return gameData, nil
	}
	if gameData.Nine == move {
		gameData.Nine = symbol
		return gameData, nil
	}
	return gameData, fmt.Errorf("not correct chose")
}
