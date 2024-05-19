package database

import (
	"database/sql"
	"fmt"
	"wasa-1967862/service/structures"
)

// Query to search a user with the unique nickname (also partially)
func (db *appdbimpl) SearchUserFromNick(nickname string, idUser int) ([]structures.User, error) {

	// Search a users
	rows, err := db.c.Query("SELECT * FROM users WHERE nickname LIKE ? AND id_user != ?", nickname+"%", idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The slice for the user
	var users []structures.User

	for rows.Next() {
		var user structures.User
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

// Query to search a user with the unique nickname and obtain the id
func (db *appdbimpl) ObtainIDFromNick(nickname string) (float64, error) {
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

// Query to search a user with the unique nickname
func (db *appdbimpl) SearchUser(nickname string) (int, error) {
	var id int
	row := db.c.QueryRow(`SELECT id_user  FROM users   WHERE nickname = ?`, nickname)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	fmt.Println("id dopo la ricerca dell utante nel db e,", id)
	return id, nil
}

func (db *appdbimpl) ExistsUser(id int) (bool, error) {
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

// Query to search a user from id to obtain the nickname
func (db *appdbimpl) GetNickname(id int) (string, error) {
	var nickname string
	row := db.c.QueryRow(`SELECT nickname  FROM users   WHERE id_user = ?`, id)
	err := row.Scan(&nickname)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return nickname, nil
}
