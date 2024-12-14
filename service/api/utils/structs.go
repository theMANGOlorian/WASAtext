package utils

// ----- Bodyrequest and BodyResponses ------

/*DoLogin structs */
type DoLoginRequestBody struct {
	Username string `json:"username"`
}

type DoLoginResponseBody struct {
	Identifier string `json:"identifier"`
	Username string `json:"username"`
	PhotoCode string `json:"photoCode"`
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

/*addToGroup structs*/
type AddToGroupRequestBody struct {
	UserId string `json:"userId"`
}
type AddToGroupResponseBody struct {
	UserId string `json:"userId"`
}

/*setMyName structs*/
type SetGroupNameRequestBody struct {
	Name string `json:"name"`
}
type SetGroupNameResponseBody struct {
	Name string `json:"name"`
}

/*setGroupPhoto structs*/
type SetGroupPhotoResponseBody struct {
	ImageCode string `json:"imageCode"`
}

/*setGroupPhoto structs*/
type SendMessageRequestBody struct {
	BodyMessage string `json:"bodyMessage"`
	ReplyTo     string `json:"replyTo"`
}

type Reactions struct {
	Emoji    string `json:"emoji"`
	Username string `json:"username"`
}

type SendMessageResponseBody struct {
	MessageId   string      `json:"messageId"`
	SenderId    string      `json:"senderId"`
	Username    string      `json:"username"`
	Text        string      `json:"text"`
	Image       string      `json:"image"`
	ReplyTo     string      `json:"replyTo"`
	Timestamp   string      `json:"timestamp"`
	Status      string      `json:"status"`
	TypeContent string      `json:"typeContent"`
	Reactions   []Reactions `json:"reactions"`
}

type Message struct {
	MessageId   string      `json:"messageId"`
	SenderId    string      `json:"senderId"`
	Username    string      `json:"username"`
	Text        string      `json:"text"`
	Image       string      `json:"image"`
	ReplyTo     string      `json:"replyTo"`
	Timestamp   string      `json:"timestamp"`
	Status      string      `json:"status"`
	TypeContent string      `json:"typeContent"`
	Reactions   []Reactions `json:"reactions"`
}

type GetConversationResponseBody struct {
	Messages   []Message `json:"messages"`
	NextCursor string    `json:"nextCursor"`
}

type ForwardMessageRequestBody struct {
	ConversationId string `json:"conversationId"`
}

type CommentMessageRequestBody struct {
	Reaction string `json:"reaction"`
}
