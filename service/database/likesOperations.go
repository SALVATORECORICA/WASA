package database

// Query to insert a like

func (db *appdbimpl) PutLike(id_photo int, id_user int) error {
	_, err := db.c.Exec("INSERT INTO likes (id_photo, id_user) VALUES (?,?)", id_photo, id_user)
	if err != nil {
		return err
	}
	return nil
}

// Check of the following exists
func (db *appdbimpl) ExistsLike(id_user int, photo_id int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM likes WHERE id_user = ? AND photo_id = ?)", id_user, photo_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) DeleteLike(id_user int, photo_id int) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE id_user = ? AND photo_id = ?", id_user, photo_id)
	if err != nil {
		return err
	}
	return nil
}
