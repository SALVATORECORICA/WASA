package database

// Query to insert a new comment on the db
func (db *appdbimpl) PostComment(id_photo int, id_user int, comment string) error {
	err := db.c.Exec("INSERT INTO comments (id_photo, id_user, comment) VALUES (?,?,?)", id_photo, id_user, comment)
	if err != nil {
		return err
	}
	return nil
}

// Query to check of the comment exist in the db
func (db *appdbimpl) ExistsComment(comment_id int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT * FROM comments WHERE comment_id = ?)", comment_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Query to check of the user is allowed to delete the comment
func (db *appdbimpl) OwnerComment(commentId int, userId int) (bool, error) {
	var idPhotoOwnerComment string
	err := db.c.Query("SELECT id_user FROM comments WHERE id_comment = ?)", commentId).Scan(&idPhotoOwnerComment)
	if err != nil {
		return false, err
	}
	if userId != idPhotoOwnerComment {
		return false, nil
	}
	return true, nil
}

func (db *appdbimpl) OwnerPhoto(userId int, photoId int) (bool, error) {

	var idPhotoOwnerPhoto string
	err := db.c.Query("SELECT id_user FROM photos WHERE id_photo = ?)", photoId).Scan(&idPhotoOwnerComment)
	if err != nil {
		return false, err
	}
	if userId != idPhotoOwnerPhoto {
		return false, nil
	}
	return true, nil
}

func (db *appdbimpl) DeleteComment(idComment int) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE id_comment = ?)", idComment)
	if err != nil {
		return err
	}
	return nil
}
