package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	conversationId := ps.ByName("conversationId")

	var requestBody utils.SendMessageRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	}
	if !utils.CheckIdentifier(conversationId) {
		ctx.Logger.WithError(err).Error("Error: conversation ID not valid")
		http.Error(w, "Error: conversation ID not valid", http.StatusBadRequest)
		return
	}
	if !utils.CheckIdentifier(requestBody.SenderId) || !(requestBody.SenderId == auth) {
		ctx.Logger.WithError(err).Error("Error: sender ID not valid")
		http.Error(w, "sender ID not valid, wrong format", http.StatusBadRequest)
		return
	}
	if requestBody.ReplyTo != "" {
		fmt.Println(requestBody.ReplyTo)
		if !utils.CheckIdentifier(requestBody.ReplyTo) {
			ctx.Logger.WithError(err).Error("Error: replyId not valid")
			http.Error(w, "replyId not valid, wrong format", http.StatusBadRequest)
			return
		}
	}
	if requestBody.BodyMessage == "" {
		ctx.Logger.WithError(err).Error("Empty messages are not allowed")
		http.Error(w, "Empty messages are not allowed", http.StatusBadRequest)
		return
	}

	code, response, err := rt.db.SendMessage(requestBody.SenderId, conversationId, requestBody.BodyMessage, requestBody.ReplyTo)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("Error: user/group not found")
			http.Error(w, "Error: user/group not found", http.StatusNotFound)
			return
		}
		if code == 403 {
			ctx.Logger.WithError(err).Error("Error: Operation forbidden")
			http.Error(w, "Error: operation forbidden", http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("An error occurred")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: encoding json ")
		http.Error(w, "Error: writing response", http.StatusInternalServerError)
		return
	}

	return
}
