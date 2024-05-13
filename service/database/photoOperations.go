package database

import (
	"encoding/base64"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/Struct"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
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

	idStr := strconv.FormatInt(id, 64)
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
func (db *appdbimpl) OwnerPhotoExtractId(photoId int) (User, error) {
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
	var path string
	err := db.c.QueryRow("SELECT path FROM photos WHERE id_photo = ?", photoId).Scan(&path)
	if err != nil {
		return path, err
	}
	return path, nil
}

func (db *appdbimpl) DeletePhoto(photoId int) error {

	_, err := db.c.Exec("DELETE FROM photos WHERE id_photo = ?", photoId)
	if err != nil {
		return err
	}
	return nil
}

// Get all photos from the id

func (db *appdbimpl) GetPhotosProfileSorted(idProfileSearched int) ([]Struct.Image, error) {
	var photos []Struct.Image
	rows, err := db.c.Query("SELECT path FROM photos WHERE id_user = ? ORDER BY date DESC ", idProfileSearched)
	if err != nil {
		return photos, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo Struct.Image
		var path string
		if err := rows.Scan(&path); err != nil {
			return photos, err
		}

		// read the file
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return photos, err
		}

		// encode Photo
		encodedData := base64.StdEncoding.EncodeToString(data)
		photo.Photo_data = encodedData
		photos = append(photos, photo)
	}
	if err = rows.Err(); err != nil {
		return photos, err
	}
	return photos, nil

}

func (db *appdbimpl) GetStream(idProfile int) ([]Image, error) {
	var photos []Image
	rows, err := db.c.Query("SELECT path FROM followers f, photos p WHERE followed_id = ? AND f.follower_id = u.id_user ORDER BY date DESC ", idProfile)
	if err != nil {
		return photos, err
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return photos, err
		}

		// read the file
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return photos, err
		}

		// encode Photo
		encodedData := base64.StdEncoding.EncodeToString(data)
		photos = append(photos, encodedData)
	}
	return photos, nil
}
