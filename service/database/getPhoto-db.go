package database

func (db *appdbimpl) GetPhoto(userId string, photo string) (bool, error) {

	const query1 = `
					SELECT EXISTS (
						-- Is your own photo profile?
						SELECT 1
						FROM users
						WHERE id = ? AND photo = ?
					) 
					OR EXISTS (
						-- Is a photo of a group?
						SELECT 1
						FROM conversations c
						JOIN members m ON m.conversationId = c.id
						WHERE m.userId = ? AND c.photo = ?
					)
					OR EXISTS(
						-- Is a photo of another user?
						SELECT 1
						FROM members m1
						JOIN members m2 ON m1.conversationId = m2.conversationId
						JOIN users u1 ON u1.id = m1.userId
						JOIN users u2 ON u2.id = m2.userId
						WHERE m1.userId = ?  -- Current user
						AND u2.photo = ?  -- The photoId of the other user
					)
					OR EXISTS (
						-- Is a photo of a message?
						SELECT 1
						FROM members m
						JOIN conversations c ON m.conversationId = c.id
						JOIN messages msg ON msg.conversation = c.id
						WHERE m.userId = ? AND msg.photo = ? 
					);
					`
	var exists int
	err := db.c.QueryRow(query1, userId, photo, userId, photo, userId, photo, userId, photo).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}

}
