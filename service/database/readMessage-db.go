package database

func (db *appdbimpl) SetReadMessage(userId string, conversationId string) (int, error) {

	const query1 = `SELECT EXISTS( SELECT 1 FROM members WHERE conversationId = ? AND userId = ?)`
	const query2 = `UPDATE members SET tsLastRead = datetime('now','localtime') WHERE userId = ? AND conversationId = ?`

	var exists int
	err := db.c.QueryRow(query1, conversationId, userId).Scan(&exists)
	if err != nil {
		return 500, err
	}

	if exists == 1 {
		_, err = db.c.Exec(query2, userId, conversationId)
		if err != nil {
			return 500, err
		}
	} else {
		return 404, nil
	}

	return 204, nil

}
