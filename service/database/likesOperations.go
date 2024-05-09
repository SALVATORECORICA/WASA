package database

// the likes
type Like struct {
	IdPhoto int `json:"id_Photo"`
	IdUser  int `json:"id_User"`
}

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

func (db *appdbimpl) GetLikes(photoId int) ([]Like, int, error) {
	var likes []Like
	var totalLikes int
	rows, err := db.c.Query("SELECT id_photo, id_user FROM likes WHERE id_photo =?", photoId)
	if err != nil {
		return likes, 0, err
	}
	// defer the closing of the rows
	defer rows.Close()
	for rows.Next() {
		var like Like
		if err := rows.Scan(&like.IdPhoto, &like.IdUser); err != nil {
			return likes, 0, err
		}
		likes = append(likes, like)
		totalLikes++
	}
	return likes, totalLikes, nil
}
