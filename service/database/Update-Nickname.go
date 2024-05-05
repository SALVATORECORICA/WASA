package database

func (db *appdbimpl) PutNewNickname(nicknameNew string, idUser int) error {

	_, err := db.c.Exec("UPDATE users  SET nickname =? WHERE id_user =?", nicknameNew, idUser)
	if err != nil {
		return err
	} else {
		return nil
	}
}
