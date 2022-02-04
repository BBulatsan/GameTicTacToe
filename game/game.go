package game

import (
	"math/rand"
	"strconv"
	"time"

	"GameTicTacToe/dbs"
)

func CheckWin(data *dbs.GameData) (bool, string, error) {
	//по горизонтале
	if data.One == data.Two && data.Two == data.Three {
		return true, data.One, nil
	}
	if data.Four == data.Five && data.Five == data.Six {
		return true, data.Four, nil
	}
	if data.Seven == data.Eight && data.Eight == data.Nine {
		return true, data.Seven, nil
	}
	//по вертикали
	if data.One == data.Four && data.Four == data.Seven {
		return true, data.One, nil
	}
	if data.Two == data.Five && data.Five == data.Eight {
		return true, data.Two, nil
	}
	if data.Three == data.Six && data.Six == data.Nine {
		return true, data.Three, nil
	}
	//по диагонали
	if data.One == data.Five && data.Five == data.Nine {
		return true, data.One, nil
	}
	if data.Three == data.Five && data.Five == data.Seven {
		return true, data.Three, nil
	}

	return false, "", nil
}

func GenUsedCk() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}
