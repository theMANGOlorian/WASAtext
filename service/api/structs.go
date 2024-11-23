package api

// ----- Bodyrequest and BodyResponses ------

/*DoLogin structs */
type doLoginRequestBody struct {
    Username string `json:"username"`
}

type doLoginResponseBody struct {
	Identifier string `json:"identifier"`
}



