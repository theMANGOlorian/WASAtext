package database

func (db *appdbimpl) DeleteMessage(userId string, messageId string) (int, error) {

	const query1 = `DELETE FROM messages WHERE id = ? AND sender = ?`

	_, err := db.c.Exec(query1, messageId, userId)
	if err != nil {
		return 500, err
	}

	return 204, nil

}
