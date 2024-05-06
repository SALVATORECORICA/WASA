package database

// Query to insert a new comment on the db
func (db *appdbimpl) PostComment(id_photo int, id_user int, comment string) error {
	err := db.c.Exec("INSERT INTO comments (id_photo, id_user, comment) VALUES (?,?,?)", id_photo, id_user, comment)
	if err != nil {
		return err
	}
	return nil
}
