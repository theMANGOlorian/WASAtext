package database

func (db *appdbimpl) GetUsersList(userId string) (*[]string, error) {

	const query1 = `SELECT username FROM users WHERE id != ?`

	rows, err := db.c.Query(query1, userId)
	if err != nil {
		return nil, err
	}

	var names []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &names, nil

}
