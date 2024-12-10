package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (db *appdbimpl) GetConversation(userId string, conversationId string, limit int, cursor string) (*utils.GetConversationResponseBody, int, error) {

	const query1 = `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? LIMIT 1`
	const query3 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query4 = `
					SELECT m.id, m.sender, u.username, m.type, m.text, m.photo, m.reply, m.status, timestamp
					FROM messages m
					JOIN users u ON u.id = m.sender
					JOIN members mb ON mb.conversationId = m.conversation
					WHERE mb.userId = ? 
						AND mb.conversationId = ? 
						AND m.timestamp < ? 
						AND m.timestamp >= mb.joinDate
					ORDER BY timestamp DESC
					LIMIT ?`

	const query5 = `
					SELECT m.id, m.sender, u.username, m.type, m.text, m.photo, m.reply, m.status, timestamp
					FROM messages m
					JOIN users u ON u.id = m.sender
					JOIN members mb ON mb.conversationId = m.conversation
					WHERE mb.userId = ? 
						AND mb.conversationId = ? 
						AND m.timestamp >= mb.joinDate
					ORDER BY timestamp DESC
					LIMIT ?`

	const query6 = `SELECT u.username, r.emoji FROM reactions r, users u WHERE r.messageId = ? AND u.id = r.owner`
	const query7 = `SELECT timestamp FROM messages WHERE id = ?`

	var exists int
	err := db.c.QueryRow(query1, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 404, fmt.Errorf("user not exists")
		}
		return nil, 500, err
	}

	err = db.c.QueryRow(query2, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 404, fmt.Errorf("conversation not exists")
		}
		return nil, 500, err
	}
	err = db.c.QueryRow(query3, userId, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 404, fmt.Errorf("user is not a member")
		}
		return nil, 500, err
	}

	var rows *sql.Rows // dichiarazione della variabile rows prima del blocco if/else

	if cursor == "" {
		rows, err = db.c.Query(query5, userId, conversationId, limit)
		if err != nil && !(errors.Is(err, sql.ErrNoRows)) {
			return nil, 500, err
		}
	} else {
		var timeCursor string
		err = db.c.QueryRow(query7, cursor).Scan(&timeCursor)
		if err != nil {
			return nil, 500, err
		}

		rows, err = db.c.Query(query4, userId, conversationId, timeCursor, limit)
		if err != nil && !(errors.Is(err, sql.ErrNoRows)) {
			return nil, 500, err
		}

	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("An error occurred while closing rows: %v", err)
		}
	}(rows)

	var nextCursor string
	var messages []utils.Message

	for rows.Next() {
		var text sql.NullString
		var image sql.NullString
		var replyTo sql.NullString
		var msg utils.Message
		if err := rows.Scan(&msg.MessageId, &msg.SenderId, &msg.Username, &msg.TypeContent, &text, &image, &replyTo, &msg.Status, &msg.Timestamp); err != nil {
			return nil, 500, err
		}
		if text.Valid {
			msg.Text = text.String
		} else {
			msg.Text = "" // O un valore di default, se necessario
		}
		if image.Valid {
			msg.Image = image.String
		} else {
			msg.Image = ""
		}
		if replyTo.Valid {
			msg.ReplyTo = replyTo.String
		} else {
			msg.ReplyTo = ""
		}

		reactionsRows, err := db.c.Query(query6, msg.MessageId)
		if err != nil && !(errors.Is(err, sql.ErrNoRows)) {
			return nil, 500, err
		}

		var reactions []utils.Reactions

		for reactionsRows.Next() {
			var reaction utils.Reactions
			if err := reactionsRows.Scan(&reaction.Username, &reaction.Emoji); err != nil {
				return nil, 500, err
			}
			reactions = append(reactions, reaction)
		}
		msg.Reactions = reactions

		messages = append(messages, msg)
		nextCursor = msg.MessageId
	}

	if err := rows.Err(); err != nil {
		return nil, 500, err
	}

	response := utils.GetConversationResponseBody{
		Messages:   messages,
		NextCursor: nextCursor,
	}

	return &response, 200, nil
}

func (db *appdbimpl) SetRecvMessage(userId string, conversationId string) error {

	// questa query chiama un trigger
	const query1 = `UPDATE members SET tsLastRecv = datetime('now','localtime') WHERE userId = ? AND conversationId = ?`
	_, err := db.c.Exec(query1, userId, conversationId)
	if err != nil {
		return err
	}
	return nil
}
