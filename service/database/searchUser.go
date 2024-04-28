package database

// Query to search a user with the unique nickname
func (db *appdbimpl) SearchUser(nickname string) (int, error) {
	var id int
	row := db.c.QueryRow(`SELECT id  FROM users   WHERE nickname = ?`, nickname)
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}
