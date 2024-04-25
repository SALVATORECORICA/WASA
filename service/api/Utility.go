package api

// Utility function
func isValidID(string id) bool {
	if len(id) >= 3 && len(id) <= 16 {
		return true
	} else {
		return false
	}
}

// Utility Query

// Query to search a user with the unique nickname
func (db *appdbimpl) searchUser(string nickname) (int, error) {
	var id int
	row := db.c.QueryRow(`SELECT id  FROM users   WHERE nickname = ?`, nickname)
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// Query to insert a user
