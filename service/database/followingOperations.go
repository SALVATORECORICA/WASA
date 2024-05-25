package database

import "database/sql"

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
