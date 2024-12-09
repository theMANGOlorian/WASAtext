package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		http.Error(w, "message ID not valid", http.StatusBadRequest)
		return
	}

	var request utils.ForwardMessageRequestBody
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("Decoding json")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.CheckIdentifier(request.ConversationId) {
		ctx.Logger.Error("Invalid ID")
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	code, response, err := rt.db.ForwardMessage(auth, messageId, request.ConversationId)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("Error: user/group not found")
			http.Error(w, "resource not found", http.StatusNotFound)
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
		ctx.Logger.WithError(err).Error("Encoding json")
		http.Error(w, "Encding json", http.StatusInternalServerError)
		return
	}
	
}
