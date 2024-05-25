package database

import "database/sql"

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
func (db *appdbimpl) DeleteLike(idUser int, photoId int) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE id_user = ? AND id_photo = ?", idUser, photoId)
	if err != nil {
		return err
	}
	return nil
}
