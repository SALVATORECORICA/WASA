package database

import "database/sql"

// Query to search a user with the unique nickname
func (db *appdbimpl) SearchUser(nickname string) (float64, error) {
	var id float64
	row := db.c.QueryRow(`SELECT id_user  FROM users   WHERE nickname = ?`, nickname)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}
