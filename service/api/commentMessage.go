package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}
	messageId := ps.ByName("messageId")

	if !utils.CheckIdentifier(messageId) {
		ctx.Logger.WithError(err).Error("Error: messageId not valid")
		http.Error(w, "Error: messageId not valid", http.StatusBadRequest)
		return
	}

	var request utils.CommentMessageRequestBody
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("Decoding Json")
		http.Error(w, "Error decoding Json", http.StatusBadRequest)
		return
	}

	if !utils.CheckReactions(request.Reaction) {
		ctx.Logger.WithError(err).Error("Reaction not allowed")
		http.Error(w, "Reaction not allowed", http.StatusBadRequest)
		return
	}

	_, err = rt.db.CommentMessage(auth, messageId, request.Reaction)
	if err != nil {
		ctx.Logger.WithError(err).Error("an error occurred")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
