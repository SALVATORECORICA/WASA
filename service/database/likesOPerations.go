package database

import (
	"database/sql"
	"wasa-1967862/service/structures"
)

// Query to insert a like

func (db *appdbimpl) PutLike(idPhoto int, idUser int) error {
	_, err := db.c.Exec("INSERT INTO likes (id_photo, id_user) VALUES (?,?)", idPhoto, idUser)
	if err != nil {
		return err
	}
	return nil
}

// Check of the like exists
func (db *appdbimpl) ExistsLike(idUser int, photoId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM likes WHERE id_user = ? AND id_photo = ?)", idUser, photoId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// Delete the like
func (db *appdbimpl) DeleteLike(idUser int, photoId int) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE id_user = ? AND id_photo = ?", idUser, photoId)
	if err != nil {
		return err
	}
	return nil
}

// Obtain all User who gives a like
func (db *appdbimpl) GetLikes(photoId int) ([]structures.User, int, error) {
	var users []structures.User
	var totalLikes int
	rows, err := db.c.Query("SELECT u.id_user, u.nickname FROM likes l JOIN users u ON l.id_user = u.id_user WHERE l.id_photo =?", photoId)
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
		totalLikes++
	}
	return users, totalLikes, nil
}

// Delete all the likes of a photo
func (db *appdbimpl) DeleteLikePhoto(idPhoto int) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE id_photo = ?)", idPhoto)
	if err != nil {
		return err
	}
	return nil
}
