package database

func (db *appdbimpl) CheckBan(u []User, idUser int) ([]User, error) {
	var updatedUsers []User
	var banner []int
	rows, err := db.c.Query(`SELECT banner_id  FROM banned_users   WHERE banned_id = ?`, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&banner); err != nil {
			return nil, err
		}
	}
	for _, user := range u {
		insert := true
		for _, b := range banner {
			if user.Id == b {
				insert = false
				break
			}
		}
		if insert {
			updatedUsers = append(updatedUsers, user)
		}
	}
	return updatedUsers, nil
}
