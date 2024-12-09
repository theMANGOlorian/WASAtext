package database

func (db *appdbimpl) ReadMessage(userId string, messageId string) (int, error) {

	const query1 = `SELECT EXISTS( SELECT 1 FROM messages WHERE messageId = ? AND userId = ?)`
	const query2 = `UPDATE users_msg SET read = 1 WHERE userId = ? AND messageId = ?`

	var exists int
	err := db.c.QueryRow(query1, messageId, userId).Scan(&exists)
	if err != nil {
		return 500, err
	}

	if exists == 1 {
		_, err = db.c.Exec(query2, userId, messageId)
		if err != nil {
			return 500, err
		}
	} else {
		return 404, nil
	}

	return 200, nil

}
