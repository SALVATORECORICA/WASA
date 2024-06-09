package database

import (
	"database/sql"
	"errors"
	"wasa-1967862/service/structures"
)

// Get the list of users whose profiles can be viewed from a list of users
// that match the search performed

func (db *appdbimpl) CheckBan(u []structures.User, idUser int) ([]structures.User, error) {
	var updatedUsers []structures.User
	var banner []int
	rows, err := db.c.Query(`SELECT banner_id  FROM banned_users   WHERE banned_id = ?`, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&banner); err != nil {
			return nil, err
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, user := range u {
		insert := true
		for _, b := range banner {
			if user.Id == b {
				insert = false
				break
			}
		}
		if insert {
			updatedUsers = append(updatedUsers, user)
		}
	}
	return updatedUsers, nil
}

// To put a new Ban
func (db *appdbimpl) PutNewBan(id_banner int, id_banned int) error {
	_, err := db.c.Exec("INSERT INTO banned_users (banner_id, banned_id) VALUES (?,?)", id_banner, id_banned)
	return err
}

// Check of a ban exists
func (db *appdbimpl) ExistsBan(id_banner int, id_banned int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM banned_users WHERE banner_id = ? AND banned_id = ?)", id_banner, id_banned).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) DeleteBan(id_banner int, id_banned int) error {
	_, err := db.c.Exec("DELETE FROM banned_users WHERE banner_id = ? AND banned_id = ?", id_banner, id_banned)
	if err != nil {
		return err
	}
	return nil
}
