package dbs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//CREATE TABLE IF NOT EXISTS games (id INTEGER PRIMARY KEY, status message_text, count_players SMALLINT, player_x_id INTEGER, player_o_id INTEGER);
//CREATE TABLE IF NOT EXISTS moves (game_id INTEGER, map_of_moves json, who_move message_text, FOREIGN KEY (game_id)  REFERENCES games (id) ON DELETE CASCADE);
//CREATE TABLE IF NOT EXISTS users ( id  integer constraint user_pk primary key, name text );
const (
	NewGame  = "newGame"
	Running  = "running"
	Finished = "finished"
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

type users struct {
	id   int
	name string
	ck   string
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

func (d *DbConn) CreateNewGame(ck string) (*GameData, error) {
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
	statement := "INSERT INTO games (status, count_players, player_x_id) VALUES ($1, $2, $3);"
	res, err := d.conn.Exec(statement, NewGame, 1, ck)
	if err != nil {
		return moveMap, err
	}
	id, _ := res.LastInsertId()
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

func (d *DbConn) RefreshGameData(gameId string) (*GameData, error) {
	gameData := &GameData{}
	var moveMap string
	rows, err := d.conn.Query(fmt.Sprintf("SELECT game_data FROM moves WHERE game_id=%s", gameId))
	if err != nil {
		return gameData, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&moveMap)
	if err != nil {
		return gameData, err
	}
	err = json.Unmarshal([]byte(moveMap), gameData)
	if err != nil {
		return gameData, err
	}
	return gameData, nil
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
		gameData, err = executeMove(gameData, move, "X")
		if err != nil {
			return gameData, count, nil
		}
		gameData.Who = "O"
	} else {
		gameData, err = executeMove(gameData, move, "O")
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

func (d *DbConn) SetGameStatus(gameId string, status string) error {
	statement := "UPDATE games SET status=$1 WHERE id=$2;"
	_, err := d.conn.Exec(statement, status, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) CreateUser(ck string) (int, error) {
	statement := "INSERT INTO users (ck) VALUES ($1);"
	res, err := d.conn.Exec(statement, ck)
	id, _ := res.LastInsertId()
	if err != nil {
		return int(id), err
	}
	return int(id), nil
}

func (d *DbConn) GetUserName(ck string) (string, error) {
	user := users{}
	rows, err := d.conn.Query(fmt.Sprintf("SELECT name FROM users WHERE ck=%s", ck))
	if err != nil {
		return user.name, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&user.name)
	return user.name, nil
}
func (d *DbConn) AddName(ck string, name string) error {
	statement := "UPDATE users SET name=$1 WHERE ck=$2;"
	_, err := d.conn.Exec(statement, name, ck)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) CheckGame(gameId string) bool {
	var count string
	rows, err := d.conn.Query(fmt.Sprintf("SELECT Count(*) FROM games WHERE id=%s", gameId))
	if err != nil {
		return false
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	i, _ := strconv.Atoi(count)
	if i == 1 {
		return true
	}
	return false
}

func (d *DbConn) SetPlayerId(userId string, gameId string) error {
	statement := "UPDATE games SET player_o_id=$1 WHERE id=$2;"
	_, err := d.conn.Exec(statement, userId, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) SetGameCount(gameId string, symbol string) error {
	c, err := d.GetCount(gameId)
	if err != nil {
		return err
	}
	if symbol == "-" {
		c--
	} else {
		c++
	}
	statement := "UPDATE games SET count_players=$1 WHERE id=$2;"
	_, err = d.conn.Exec(statement, c, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) GetCount(gameId string) (int, error) {
	var count string
	rows, err := d.conn.Query(fmt.Sprintf("SELECT count_players FROM games WHERE id=%s", gameId))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	c, _ := strconv.Atoi(count)
	return c, nil
}

func executeMove(gameData *GameData, move string, symbol string) (*GameData, error) {
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
