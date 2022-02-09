package dbs

type moves struct {
	GameId    int
	GameData  string
	CountMove int
}

func (d *DbConn) CreateMove(gd []byte, count int, gameId string) error {
	statement := "INSERT INTO moves (game_id, game_data, count_move) VALUES ($1, $2, $3)"
	_, err := d.conn.Exec(statement, gameId, gd, count)
	if err != nil {
		return err
	}
	return nil
}

func (d *DbConn) SetMove(gd []byte, count int, gameId string) error {
	statement := "UPDATE moves SET game_data=$1, count_move=$2 WHERE game_id=$3;"
	_, err := d.conn.Exec(statement, gd, count, gameId)
	if err != nil {
		return err
	}
	return nil
}
