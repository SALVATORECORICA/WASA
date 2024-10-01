package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
	"wasa-1967862/service/structures"
)

// Query to insert a new photo on the db and return the id of the new photo inserted
func (db *appdbimpl) PostNewPhoto(id_user int, path string, timestamp time.Time) (int, string, error) {
	result, err := db.c.Exec("INSERT INTO photos (id_user, uploadDate, path) VALUES (?,?,?)", id_user, timestamp, path)
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
	err := db.c.QueryRow("SELECT id_user FROM photos WHERE  id_photo = ?", photoId).Scan(&user.Id)

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
	var uploadDate time.Time
	err := db.c.QueryRow("SELECT uploadDate FROM photos WHERE id_photo = ?", photoId).Scan(&uploadDate)
	if err != nil {
		return uploadDate, err
	}
	return uploadDate, nil
}

// Obtain the path of the Photo
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

func (db *appdbimpl) GetPhotosProfileSorted(idProfileSearched int) ([]structures.Photo, error) {
	var photos []structures.Photo
	rows, err := db.c.Query("SELECT path FROM photos WHERE id_user = ? ORDER BY uploadDate DESC ", idProfileSearched)
	if err != nil {
		return photos, err
	}
	defer rows.Close()
	for rows.Next() {
		var idPhoto int
		if err := rows.Scan(&idPhoto); err != nil {
			return photos, err
		}
		photo, err := db.GetPhotoComplete(idPhoto, idProfileSearched)
		if err != nil {
			return photos, err
		}
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		return photos, err
	}
	return photos, nil
}

func (db *appdbimpl) GetStream(idProfile int) ([]structures.Photo, error) {
	var photos []structures.Photo
	rows, err := db.c.Query("SELECT id_photo FROM followers f, photos p WHERE follower_id = ? AND f.followed_id = p.id_user ORDER BY uploadDate DESC ", idProfile)
	if err != nil {
		return photos, err
	}
	defer rows.Close()
	for rows.Next() {
		var idPhoto int
		if err := rows.Scan(&idPhoto); err != nil {
			return photos, err
		}
		photo, err := db.GetPhotoComplete(idPhoto, idProfile)
		if err != nil {
			return photos, err
		}
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		return photos, err
	}
	return photos, nil
}

func (db *appdbimpl) GetPhotoComplete(photoId int, idUser int) (structures.Photo, error) {

	// Declare the photo to return
	var photo structures.Photo

	// Obtain the owner of the photo

	owner, err := db.OwnerPhotoFromIdPhoto(photoId)
	if err != nil {
		return photo, err
	}

	// obtain the likes

	usersLikes, nLikes, err := db.GetLikes(photoId)
	if err != nil {
		return photo, err
	}
	// obtain the comments
	comments, err := db.CommentsPhoto(photoId)
	if err != nil {
		return photo, err
	}

	// Obtain the date
	date, err := db.GetPhotoDate(photoId)
	if err != nil {
		return photo, err
	}
	// Obtain the image

	// Obtain the path where we can find the photo
	path, err := db.GetPhotoPath(photoId)
	if err != nil {
		return photo, err
	}

	// Take the photo

	photoData, err := ioutil.ReadFile(path)
	if err != nil {
		return photo, err
	}

	// Check if the user liked the photo
	existsLike, err := db.ExistsLike(idUser, photoId)
	if err != nil {
		return photo, err
	}

	publicPath := "/static/" + filepath.Base(path)

	// Now we are ready to send the

	photo.PhotoId = photoId
	photo.Owner = owner
	photo.Date = date
	photo.Likes = usersLikes
	photo.NLikes = nLikes
	photo.PhotoData = photoData
	photo.Comments = comments
	photo.Liked = existsLike
	photo.Path = publicPath
	fmt.Println("eccolo", path)

	return photo, nil
}
