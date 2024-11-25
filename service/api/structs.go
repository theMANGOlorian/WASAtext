package api

// ----- Bodyrequest and BodyResponses ------

/*DoLogin structs */
type doLoginRequestBody struct {
	Username string `json:"username"`
}

type doLoginResponseBody struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type setMyUserNameRequestBody struct {
	Username string `json:"username"`
}

type setMyUserNameResponseBody struct {
	Username string `json:"username"`
}
