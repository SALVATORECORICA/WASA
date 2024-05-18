package database

import (
	"path/filepath"
	"strconv"
	"time"
	"wasa-1967862/service/structures"
)

// Query to insert a new photo on the db and return the id of the new photo inserted
func (db *appdbimpl) PostNewPhoto(nickname string, path string, timestamp time.Time) (int, string, error) {
	result, err := db.c.Exec("INSERT INTO photos (id_user, date, path) VALUES (?,?,?)", nickname, timestamp, path)
	if err != nil {
		return -1, "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, "", err
	}
	idStr := strconv.FormatInt(id, 10)
	idInt := int(id)
	// create the complete path of the photo
	completePath := filepath.Join(path, idStr+".jpg")
	// Update the path
	_, err = db.c.Exec("UPDATE photos SET path = ? WHERE id_photo= ?", completePath, id)
	if err != nil {
		return -1, "", err
	}
	return idInt, completePath, nil
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
func (db *appdbimpl) OwnerPhotoFromIdPhoto(photoId int) (structures.User, error) {
	var user structures.User
	err := db.c.QueryRow("SELECT id_user FROM photos WHERE  id_photo = ?)", photoId).Scan(&user.Id)
	if err != nil {
		return user, err
	}
	err = db.c.QueryRow("SELECT nickname FROM users WHERE id_user =?", user.Id).Scan(&user.Nickname)
	if err != nil {
		return user, err
	}
	return user, nil
}
