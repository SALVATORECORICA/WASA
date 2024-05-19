package database

// Query to insert a user

func (db *appdbimpl) PutNewUser(nickname string) (int, error) {
	_, err := db.c.Exec("INSERT INTO users (nickname) VALUES (?)", nickname)
	if err != nil {
		return 0, err
	}
	id, err := db.SearchUser(nickname)
	return id, err
}
