package database

// Query to search a user with the unique nickname (also partially)
func (db *appdbimpl) SearchUserFromNick(nickname string) ([]User, error) {

	// Search a users
	rows, err := db.c.Query("SELECT * FROM users WHERE nickname LIKE ?", nickname+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The slice for the user
	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Nickname)
		if err != nil {
			return nil, err
		}
		// Add the user to the slice
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// the Profile of the user
type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}
