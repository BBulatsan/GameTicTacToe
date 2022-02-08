package dbs

type users struct {
	id   int
	name string
	ck   string
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
	statement := "SELECT name FROM users WHERE ck=$1"
	rows := d.conn.QueryRow(statement, ck)
	err := rows.Scan(&user.name)
	if err != nil {
		return user.name, err
	}
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
