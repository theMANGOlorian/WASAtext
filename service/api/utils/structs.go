package utils

// ----- Bodyrequest and BodyResponses ------

/*DoLogin structs */
type DoLoginRequestBody struct {
	Username string `json:"username"`
}

type DoLoginResponseBody struct {
	Identifier string `json:"identifier"`
}

/*SetMyUserName structs*/
type SetMyUserNameRequestBody struct {
	Username string `json:"username"`
}

type SetMyUserNameResponseBody struct {
	Username string `json:"username"`
}

/*SetMyPhoto structs*/
type SetMyPhotoResponseBody struct {
	ImageCode string `json:"imageCode"`
}

/*StartConversation structs*/
type StartConversationRequestBody struct {
	Name             string `json:"name"`
	TypeConversation string `json:"typeConversation"`
}
type StartConversationResponseBody struct {
	Identifier string `json:"identifier"`
}

/*GetConversation structs*/
type GetConversationsResponseBody struct {
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	ConversationId       string `json:"conversationId"`
	TypeConversation     string `json:"conversationType"`
	ConversationName     string `json:"conversationName"`
	PhotoProfileCode     string `json:"photoProfileCode"`
	LastMessageTimeStamp string `json:"lastMessageTimeStamp"`
	LastMessagePreview   string `json:"lastMessagePreview"`
}
