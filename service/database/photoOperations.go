package database

import "time"

// Query to insert a new photo on the db and return the id of the new photo inserted
func (db *appdbimpl) PostNewPhoto(nickname string, complete_path string, timestamp time.Time) (int, error) {
	result, err := db.c.Exec("INSERT INTO photos (id_user, date, path) VALUES (?,?,?)", nickname, timestamp, complete_path)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}
	return id, nil
}

// Query to check if the photo Exists
func (db *appdbimpl) ExistsPhoto(photoId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM photos WHERE  id_photo = ?)", photoId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Extract the id of the owner
func (db *appdbimpl) OwnerPhoto(photoId int) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id_user FROM photos WHERE  id_photo = ?)", photoId).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
