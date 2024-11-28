package api

import (
	"WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (rt *_router) startConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		ctx.Logger.Error("Error: Missing Authorization header")
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		ctx.Logger.Error("Error: Invalid Authorization header")
		http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
		return
	}
	auth := strings.TrimPrefix(authHeader, bearerPrefix)

	// checking permission
	id := ps.ByName("userId")
	if !isAuthorized(id, auth) {
		ctx.Logger.Error("Error: User not Authorized")
		http.Error(w, "Operation Forbidden", http.StatusForbidden)
		return
	}

	if r.ContentLength == 0 || r.Header.Get("Content-Type") != "application/json" {
		ctx.Logger.Error("Error: Content-Type not supported")
		http.Error(w, "Content-Type not valid", http.StatusBadRequest)
		return
	}

	var requestBody startConversationRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	}

	var conversationId string
	if requestBody.TypeConversation == "group" {
		conversationId, err = rt.db.StartConversationGroup(id, requestBody.Name)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error starting conversation")
			http.Error(w, "Cannot start conversation", http.StatusInternalServerError)
			return
		}
	} else {
		if requestBody.TypeConversation == "private" {
			conversationId, err = rt.db.StartConversationPrivate(id, requestBody.Name)
			if err != nil {

				if err.Error() == "EXISTS" {
					ctx.Logger.WithError(err).Error("Error starting conversation")
					http.Error(w, "it already exists", http.StatusInternalServerError)
					return
				}
				if err.Error() == "NOT FOUND" {
					ctx.Logger.WithError(err).Error("Error starting conversation")
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}
				ctx.Logger.WithError(err).Error("Error starting conversation")
				http.Error(w, "Cannot start conversation", http.StatusInternalServerError)
				return
			}
		} else {
			ctx.Logger.WithError(err).Error("Error: typeConversation value not valid ")
			http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
			return
		}
	}

	// Generating response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var responseBody startConversationResponseBody
	responseBody.Identifier = conversationId
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Encoding JSON ")
		http.Error(w, "Cannot encode JSON", http.StatusInternalServerError)
	}

	return

}
