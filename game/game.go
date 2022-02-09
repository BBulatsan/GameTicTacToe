package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"GameTicTacToe/dbs"
)

func CheckWin(data *dbs.GameData) (bool, string, error) {
	// horizontally
	if data.One == data.Two && data.Two == data.Three {
		return true, data.One, nil
	}
	if data.Four == data.Five && data.Five == data.Six {
		return true, data.Four, nil
	}
	if data.Seven == data.Eight && data.Eight == data.Nine {
		return true, data.Seven, nil
	}
	// vertically
	if data.One == data.Four && data.Four == data.Seven {
		return true, data.One, nil
	}
	if data.Two == data.Five && data.Five == data.Eight {
		return true, data.Two, nil
	}
	if data.Three == data.Six && data.Six == data.Nine {
		return true, data.Three, nil
	}
	// diagonally
	if data.One == data.Five && data.Five == data.Nine {
		return true, data.One, nil
	}
	if data.Three == data.Five && data.Five == data.Seven {
		return true, data.Three, nil
	}

	return false, "", nil
}

func MakeMove(d dbs.DbConn, move string, gameId string) (*dbs.GameData, bool, int, error) {
	var unique bool
	gameData, count, err := d.GetGameData(gameId)
	if err != nil {
		return gameData, unique, count, err
	}

	if gameData.Symbol == dbs.X {
		gameData, err = executeMove(gameData, move, dbs.X)
		if err != nil {
			return gameData, unique, count, nil
		}
		gameData.Symbol = dbs.O
	} else {
		gameData, err = executeMove(gameData, move, dbs.O)
		if err != nil {
			return gameData, unique, count, nil
		}
		gameData.Symbol = dbs.X
	}

	gd, err := json.Marshal(gameData)
	if err != nil {
		return gameData, unique, count, err
	}
	count++
	err = d.SetMove(gd, count, gameId)
	if err != nil {
		return gameData, unique, count - 1, err
	}
	unique = true
	return gameData, unique, count, nil
}

func executeMove(gameData *dbs.GameData, move string, symbol string) (*dbs.GameData, error) {
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

func GenUsedCk() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}
