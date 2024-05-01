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

func (db *appdbimpl) SearchUserID(id int) (bool, error) {
	var idUser int
	row := db.c.QueryRow(`SELECT id_user  FROM users   WHERE id_user = ?`, id)
	err := row.Scan(&idUser)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
