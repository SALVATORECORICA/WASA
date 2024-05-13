package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/Struct"
)

// Query to insert a new comment on the db
func (db *appdbimpl) PostComment(id_photo int, id_user int, comment string) error {
	_, err := db.c.Exec("INSERT INTO comments (id_photo, id_user, comment) VALUES (?,?,?)", id_photo, id_user, comment)
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
	var idPhotoOwnerComment int
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

	var idPhotoOwner int
	err := db.c.QueryRow("SELECT id_user FROM photos WHERE id_photo = ?)", photoId).Scan(&idPhotoOwner)
	if err != nil {
		return false, err
	}
	if userId != idPhotoOwner {
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

// Obtain all comments of a Photo

func (db *appdbimpl) CommentsPhoto(photoId int) ([]Struct.Comments, error) {
	var comments []Struct.Comments
	rows, err := db.c.Query("SELECT id_comment, id_user, comment FROM comments WHERE id_photo=?", photoId)
	if err != nil {
		return comments, err
	}
	// defer the closing of the rows
	defer rows.Close()
	for rows.Next() {
		var comment Struct.Comments
		if err := rows.Scan(&comment.Comment_id, &comment.User.Id, &comment.Comment); err != nil {
			return comments, err
		}
		comments = append(comments, comment)

	}
	return comments, nil
}

// Delete the comments of a photo
func (db *appdbimpl) DeleteCommentPhoto(idPhoto int) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE id_photo = ?)", idPhoto)
	if err != nil {
		return err
	}
	return nil
}
