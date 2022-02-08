package dbs

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//CREATE TABLE IF NOT EXISTS games (id INTEGER PRIMARY KEY, status message_text, player_x_id INTEGER, player_o_id INTEGER);
//CREATE TABLE IF NOT EXISTS moves (game_id INTEGER, map_of_moves json, who_move message_text, FOREIGN KEY (game_id)  REFERENCES games (id) ON DELETE CASCADE);
//CREATE TABLE IF NOT EXISTS users ( id  integer constraint user_pk primary key, name text );

type DbConn struct {
	conn *sql.DB
}

func (d *DbConn) InitDb() error {
	db, err := sql.Open("sqlite3", "./dbs/store/store.db")
	if err != nil {
		return err
	}
	d.conn = db
	return nil
}
