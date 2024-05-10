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
func (db *appdbimpl) OwnerPhoto(photoId int) (User, error) {
	var user User
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

// Obtain the date of the Photo

func (db *appdbimpl) GetPhotoDate(photoId int) (time.Time, error) {
	var time time.Time
	err := db.c.QueryRow("SELECT date FROM photos WHERE id_photo = ?", photoId).Scan(&time)
	if err != nil {
		return time, err
	}
	return time, nil
}

func (db *appdbimpl) GetPhotoPath(photoId int) (string, error) {
	var path sting
	err := db.c.QueryRow("SELECT path FROM photos WHERE id_photo = ?", photoId).Scan(&path)
	if err != nil {
		return path, err
	}
	return path, nil
}

func (db *appdbimpl) DeletePhoto(photoId int) error {

	err := db.c.Exec("DELETE FROM photos WHERE id_photo = ?", photoId)
	if err != nil {
		return err
	}
	return nil
}

// Get all photos from the id

func (db *appdbimpl) GetPhotosProfile(idProfileSearched int) (string, error) {

}
