package database

func (db *appdbimpl) GetPhoto(userId string, photo string) (bool, error) {

	const query1 = `SELECT EXISTS (
						-- Query 1
						SELECT 1
						FROM users
						WHERE id = ? AND photo = ?
					) 
					OR EXISTS (
						-- Query 2
						SELECT 1
						FROM conversations c
						JOIN members m ON m.conversationId = c.id
						WHERE m.userId = ? AND c.photo = ?
					) 
					OR EXISTS (
						-- Query 3
						SELECT 1
						FROM members m
						JOIN conversations c ON m.conversationId = c.id
						JOIN messages msg ON msg.conversation = c.id
						WHERE m.userId = ? AND msg.photo = ? 
					);
					`
	var exists int
	err := db.c.QueryRow(query1, userId, photo, userId, photo, userId, photo).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}

}
