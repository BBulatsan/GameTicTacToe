package dbs

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	NewGame  = "newGame"
	Running  = "running"
	Finished = "finished"
	X        = "X"
	O        = "O"
)

type Games struct {
	Id        int
	Status    string
	PlayerXId int
	PlayerOId int
}

type GameData struct {
	Who    string
	Symbol string
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

func (d *DbConn) CreateGame() (int, error) {
	var gameId int
	statement := "INSERT INTO games (status) VALUES ($1);"
	res, err := d.conn.Exec(statement, NewGame)
	if err != nil {
		return gameId, err
	}
	id, _ := res.LastInsertId()
	gameId = int(id)
	return gameId, nil
}

//func (d *DbConn) CreateNewGame(ck string) (*GameData, error) {
//	var gameId int
//	moveMap := &GameData{
//		Symbol: "X",
//		One:    "1",
//		Two:    "2",
//		Three:  "3",
//		Four:   "4",
//		Five:   "5",
//		Six:    "6",
//		Seven:  "7",
//		Eight:  "8",
//		Nine:   "9",
//	}
//	statement := "INSERT INTO games (status, player_x_id) VALUES ($1, $2);"
//	res, err := d.conn.Exec(statement, NewGame, ck)
//	if err != nil {
//		return moveMap, err
//	}
//	id, _ := res.LastInsertId()
//	gameId = int(id)
//	moveMap.GameId = gameId
//
//	js, err := json.Marshal(moveMap)
//	if err != nil {
//		return moveMap, err
//	}
//	statement = "INSERT INTO moves (game_id, game_data, count_move) VALUES ($1, $2, $3);"
//	_, err = d.conn.Exec(statement, gameId, js, 1)
//	if err != nil {
//		return moveMap, err
//	}
//
//	return moveMap, nil
//}

func (d *DbConn) RefreshGameData(gameId string) (*GameData, error) {
	gameData := &GameData{}
	var moveMap string
	statement := "SELECT game_data FROM moves WHERE game_id=$1"
	rows := d.conn.QueryRow(statement, gameId)
	err := rows.Scan(&moveMap)
	if err != nil {
		return gameData, err
	}
	err = json.Unmarshal([]byte(moveMap), gameData)
	if err != nil {
		return gameData, err
	}
	return gameData, nil
}

func (d *DbConn) GetGameData(gameId string) (*GameData, int, error) {
	move := &moves{}
	gameData := &GameData{}
	statement := "SELECT game_data, count_move FROM moves WHERE game_id=$1"
	rows := d.conn.QueryRow(statement, gameId)
	err := rows.Scan(&move.GameData, &move.CountMove)
	if err != nil {
		return gameData, move.CountMove, err
	}
	err = json.Unmarshal([]byte(move.GameData), gameData)
	if err != nil {
		return gameData, move.CountMove, err
	}

	return gameData, move.CountMove, nil
}

func (d *DbConn) SetGameStatus(gameId string, status string) error {
	statement := "UPDATE games SET status=$1 WHERE id=$2;"
	_, err := d.conn.Exec(statement, status, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) GetGameStatus(gameId string) (string, error) {
	var status string
	statement := "SELECT status FROM games WHERE id=$1"
	rows := d.conn.QueryRow(statement, gameId)
	err := rows.Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (d *DbConn) CheckGame(gameId string) bool {
	var count string
	statement := "SELECT Count(*) FROM games WHERE id=$1"
	rows := d.conn.QueryRow(statement, gameId)
	err := rows.Scan(&count)
	if err != nil {
		return false
	}
	i, _ := strconv.Atoi(count)
	if i == 1 {
		return true
	}
	return false
}

func (d *DbConn) SetPlayerId(ck string, gameId string, symbol string) error {
	if symbol == O {
		statement := "UPDATE games SET player_o_id=$1 WHERE id=$2;"
		_, err := d.conn.Exec(statement, ck, gameId)
		if err != nil {
			return err
		}
		return nil
	} else if symbol == X {
		statement := "UPDATE games SET player_x_id=$1 WHERE id=$2;"
		_, err := d.conn.Exec(statement, ck, gameId)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("uknow symbol")
}

func (d *DbConn) GetPlayersCK(gameId string) (Games, error) {
	game := Games{}
	statement := "SELECT player_x_id, player_o_id FROM games WHERE id=$1"
	rows := d.conn.QueryRow(statement, gameId)
	err := rows.Scan(&game.PlayerXId, &game.PlayerOId)
	if err != nil {
		return game, err
	}
	return game, nil
}
