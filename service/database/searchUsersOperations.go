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

// Check of the user exists
func (db *appdbimpl) ExistsUser(id int) (bool, error) {

	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE  id_user = ?)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Query to search a user with the unique nickname (also partially)
func (db *appdbimpl) SearchUserFromNick(nickname string, idUser int) ([]User, error) {

	// Search a users
	rows, err := db.c.Query("SELECT * FROM users WHERE nickname LIKE ? AND id_user != ?", nickname+"%", idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The slice for the user
	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Nickname)
		if err != nil {
			return nil, err
		}
		// Add the user to the slice
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// the Profile of the user
type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}
