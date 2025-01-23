package database

import (
	"WASAtext/service/api/utils"
	"database/sql"
	"fmt"
	"log"
)

func (db *appdbimpl) GetConversations(userId string, r *utils.GetConversationsResponseBody) error {
	query1 := `
	SELECT 
    c.id AS conversation_id, 
    c.type AS conversation_type, 
    c.name AS conversation_name, 
    c.photo AS conversation_photo,
    COALESCE(msg.timestamp, '0000-00-00 00:00:00') AS last_message_timestamp,
    msg.type AS last_message_type,
    msg.text AS last_message_text
	FROM 
		conversations c
	JOIN 
		members m ON m.conversationId = c.id
	JOIN 
		users u ON u.id = m.userId
	LEFT JOIN 
		messages msg ON msg.conversation = c.id
		AND msg.timestamp = (SELECT MAX(timestamp) FROM messages WHERE conversation = c.id)
	WHERE 
		u.id = ? 
	ORDER BY 
		COALESCE(msg.timestamp, '0000-00-00 00:00:00') DESC;
	`
	// if there aren't message in a conversation then it will set the lastMessageTimeStamp to default data (0000-00-00 00:00:00)
	rows, err := db.c.Query(query1, userId)
	if err != nil {
		return fmt.Errorf("error while getting conversations ID : %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows: ", err)
		}
	}(rows)

	var conversations []utils.Conversation
	for rows.Next() {
		var conv utils.Conversation
		var lastMessageType sql.NullString
		var lastMessageText sql.NullString
		var conversationPhoto sql.NullString
		var lastMessageTimeStamp sql.NullString
		var name sql.NullString

		err := rows.Scan(
			&conv.ConversationId,
			&conv.TypeConversation,
			&name,
			&conversationPhoto,
			&lastMessageTimeStamp,
			&lastMessageType,
			&lastMessageText,
		)
		if err != nil {
			return fmt.Errorf("error scanning conversation row: %w", err)
		}

		if conv.TypeConversation == "group" {
			// if conversation is a group then it will use name in conversation table (in db)
			conv.ConversationName = name.String
			if conversationPhoto.Valid {
				conv.PhotoProfileCode = conversationPhoto.String
			} else {
				conv.PhotoProfileCode = ""
			}

		} else {
			// if conversation is private then it will use the name of the user he is talking to
			var friendUserPhoto sql.NullString
			query2 := `SELECT u.username, u.photo FROM members m JOIN users u ON m.userId = u.id WHERE userId != ? AND conversationId = ?`
			err := db.c.QueryRow(query2, userId, conv.ConversationId).Scan(&conv.ConversationName, &friendUserPhoto)
			if err != nil {
				return fmt.Errorf("error while getting conversation username : %w", err)
			}
			if friendUserPhoto.Valid {
				conv.PhotoProfileCode = friendUserPhoto.String
			}

		}

		if lastMessageTimeStamp.Valid {
			conv.LastMessageTimeStamp = lastMessageTimeStamp.String
		} else {
			conv.LastMessagePreview = ""
		}
		var messageType string
		if lastMessageType.Valid {
			messageType = lastMessageType.String
		}

		switch messageType {
		case "text":
			conv.LastMessagePreview = lastMessageText.String
		case "photo":
			conv.LastMessagePreview = "📸"
		default:
			conv.LastMessagePreview = ""
		}
		
		// append conversation
		conversations = append(conversations, conv)

	}

	if err := rows.Err(); err != nil {
		return err
	}

	r.Conversations = conversations
	return nil
}
