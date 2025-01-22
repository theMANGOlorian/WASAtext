package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	conversationId := ps.ByName("conversationId")

	// verify requestBody
	var requestBody utils.AddToGroupRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	}

	if !utils.CheckIdentifier(conversationId) {
		ctx.Logger.WithError(err).Error("Error: conversationId not valid")
		http.Error(w, "Error: conversationId not valid", http.StatusBadRequest)
		return
	}

	if !utils.CheckName(requestBody.Username) {
		ctx.Logger.WithError(err).Error("Error: userId not valid")
		http.Error(w, "userId not valid, wrong format", http.StatusBadRequest)
		return
	}

	code, err := rt.db.AddToGroupPermission(auth, conversationId)
	if err != nil {
		if code == 403 {
			ctx.Logger.WithError(err).Error("User doesn't have permission to add someone in the group")
			http.Error(w, "User doesn't have permission ", http.StatusForbidden)
			return
		}
		if code == 404 {
			ctx.Logger.WithError(err).Error("group or user not found")
			http.Error(w, "user or group not found", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("an error occurred")
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}
	// add users in group

	code, err = rt.db.AddToGroup(requestBody.Username, conversationId)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("user not found")
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if code == 409 {
			ctx.Logger.WithError(err).Error("user already exists")
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		ctx.Logger.WithError(err).Error("an error occurred")
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var responseBody utils.AddToGroupResponseBody
	responseBody.Username = requestBody.Username
	err = json.NewEncoder(w).Encode(&responseBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Encoding JSON ")
		http.Error(w, "encoding JSON", http.StatusInternalServerError)
		return
	}
}
