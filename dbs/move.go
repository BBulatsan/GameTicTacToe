package dbs

type moves struct {
	GameId    int
	GameData  string
	CountMove int
}

func (d *DbConn) SetMove(gd []byte, count int, gameId string) error {
	statement := "UPDATE moves SET game_data=$1, count_move=$2 WHERE game_id=$3;"
	_, err := d.conn.Exec(statement, gd, count, gameId)
	if err != nil {
		return err
	}
	return nil
}
