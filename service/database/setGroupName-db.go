package database

// controlalre user esiste
// contrllare gruppo esiste ed se è un gruppo
// contrllare se user è nel gruppo
// cambiare nome
func (db *appdbimpl) SetGroupName(userId string, groupId, name string) (int, error) {
	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? AND type = 'group' LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `UPDATE conversations SET name = ? WHERE id = ?`

	var exists int
	err := db.c.QueryRow(query1, userId).Scan(&exists)
	if err != nil {
		return 404, err
	}

	err = db.c.QueryRow(query2, groupId).Scan(&exists)
	if err != nil {
		return 404, err
	}

	err = db.c.QueryRow(query3, userId, groupId).Scan(&exists)
	if err != nil {
		return 403, err
	}

	_, err = db.c.Exec(query4, name, groupId)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
