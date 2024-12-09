package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func (db *appdbimpl) ForwardMessage(userId string, messageId string, toConversation string) (int, *utils.SendMessageResponseBody, error) {

	const query1 = `SELECT EXISTS (SELECT 1 FROM members WHERE userId = ? AND conversationId = ?)`
	const query2 = `SELECT type, text, photo FROM messages WHERE id = ?` // se non lo trova il messaggio non esiste
	const query3 = `INSERT INTO messages (id, sender, type, text, photo, conversation) VALUES (?,?,?,?,?,?)`
	const query4 = `SELECT u.username, m.timestamp FROM messages m JOIN users u ON m.sender = u.id WHERE m.id = ?`

	var exists int
	err := db.c.QueryRow(query1, userId, toConversation).Scan(&exists)
	if err != nil {
		return 500, nil, err
	}
	if exists == 0 {
		return 404, nil, fmt.Errorf("user/conversation not found or user are not a member")
	}

	var response utils.SendMessageResponseBody
	var typeContent string
	var text sql.NullString
	var photo sql.NullString

	err = db.c.QueryRow(query2, messageId).Scan(&typeContent, &text, &photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, nil, fmt.Errorf("message not found")
		} else {
			return 500, nil, err
		}
	}
	if text.Valid {
		response.Text = text.String
	} else {
		response.Text = ""
	}
	if photo.Valid {
		response.Image = photo.String
	} else {
		response.Image = ""
	}
	NewMessageId := uuid.New().String()

	for {
		_, err = db.c.Exec(query3, NewMessageId, userId, typeContent, text, photo, toConversation)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			messageId = uuid.New().String()
		} else {
			return 500, nil, err
		}
	}

	err = db.c.QueryRow(query4, NewMessageId).Scan(&response.Username, &response.Timestamp)
	if err != nil {
		return 500, nil, err
	}

	response.MessageId = NewMessageId
	response.SenderId = userId
	response.ReplyTo = ""
	response.Status = "none"
	response.TypeContent = typeContent
	response.Reactions = []utils.Reactions{}

	return 201, &response, nil
}
