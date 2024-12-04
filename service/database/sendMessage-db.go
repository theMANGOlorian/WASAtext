package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"strings"
)

// controllare se l'utente fa parte della conversazione
// aggiungere messaggio

func (db *appdbimpl) SendMessage(userId string, conversationId string, text string, replyTo ...string) (int, *utils.SendMessageResponseBody, error) {

	const query1 = `SELECT username FROM users WHERE id = ? LIMIT 1`
	const query2 = `SELECT 1 FROM conversations WHERE id = ? LIMIT 1`
	const query3 = `SELECT 1 FROM messages WHERE id = ? AND conversation = ? LIMIT 1`
	const query4 = `SELECT 1 FROM members WHERE userId = ? AND conversationId = ? LIMIT 1`
	const query5 = `INSERT INTO messages(id, type, text, conversation,reply, status) VALUES (?,'text',?,?,?,'none')`
	const query6 = `SELECT timestamp FROM messages WHERE id = ? LIMIT 1`

	if len(replyTo) > 1 {
		return 500, nil, fmt.Errorf("too many message id, you can reply JUST once")
	}
	var replyId string
	if len(replyTo) == 1 {
		replyId = replyTo[0]
	} else {
		replyId = ""
	}

	var exists int
	var username string
	err := db.c.QueryRow(query1, userId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, nil, fmt.Errorf("user not exists")
		}
		return 500, nil, err
	}

	err = db.c.QueryRow(query2, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, nil, fmt.Errorf("conversation not exists")
		}
		return 500, nil, err
	}

	if replyId != "" {
		err = db.c.QueryRow(query3, replyId, conversationId).Scan(&exists)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 404, nil, fmt.Errorf("message to reply not exists")
			}
			return 500, nil, err
		}
	}
	err = db.c.QueryRow(query4, userId, conversationId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 403, nil, fmt.Errorf("user cannot send message in this conversation because he is not a member")
		}
		return 500, nil, err
	}

	messageId := uuid.New().String()
	for {
		_, err = db.c.Exec(query5, messageId, text, conversationId, replyId)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			messageId = ksuid.New().String()
		} else {
			return 500, nil, err
		}
	}
	var ts string
	err = db.c.QueryRow(query6, messageId).Scan(&ts)
	if err != nil {
		return 500, nil, err
	}

	var response = &utils.SendMessageResponseBody{
		MessageId:   messageId,
		SenderId:    userId,
		Username:    username,
		Text:        text,
		Image:       "",
		ReplyTo:     replyId,
		Timestamp:   ts,
		Status:      "none",
		TypeContent: "text",
		Reactions:   []utils.Reactions{},
	}

	return 201, response, nil

}
