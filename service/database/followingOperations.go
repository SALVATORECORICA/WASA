package database

import (
	"database/sql"
	"wasa-1967862/service/structures"
)

// To put a new Ban
func (db *appdbimpl) PutFollowing(follower_id int, followed_id int) error {
	_, err := db.c.Exec("INSERT INTO followers (follower_id, followed_id) VALUES (?,?)", follower_id, followed_id)
	return err
}

// Check of the following exists
func (db *appdbimpl) ExistsFollowing(follower_id int, followed_id int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM followers WHERE follower_id = ? AND followed_id = ?)", follower_id, followed_id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) DeleteFollowing(follower_id int, followed_id int) error {
	_, err := db.c.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

// Obtain all followers and the number of them
func (db *appdbimpl) GetFollower(userId int) ([]structures.User, int, error) {
	var nFollowers int
	var users []structures.User
	rows, err := db.c.Query("SELECT follower_id, nickname FROM followers f, users u WHERE f.followed_id= ? AND f.follower_id = u.id_user", userId)
	if err != nil {
		return users, 0, err
	}

	// defer the closing of the rows
	defer rows.Close()
	for rows.Next() {
		var user structures.User
		if err := rows.Scan(&user.Id, &user.Nickname); err != nil {
			return users, 0, err
		}
		users = append(users, user)
		nFollowers++
	}
	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return users, 0, err
	}
	return users, nFollowers, nil
}

// Obtain all followed and the number of them
func (db *appdbimpl) GetFollowed(userId int) ([]structures.User, int, error) {

	var nFollowed int
	var users []structures.User
	rows, err := db.c.Query("SELECT followed_id, nickname FROM followers f, users u WHERE f.follower_id= ? AND f.followed_id = u.id_user", userId)
	if err != nil {
		return users, 0, err
	}

	// defer the closing of the rows
	defer rows.Close()
	for rows.Next() {
		var user structures.User
		if err := rows.Scan(&user.Id, &user.Nickname); err != nil {
			return users, 0, err
		}
		users = append(users, user)
		nFollowed++
	}
	if err := rows.Err(); err != nil {
		return users, 0, err
	}
	return users, nFollowed, nil
}
