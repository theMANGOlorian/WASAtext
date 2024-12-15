package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")

	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	conversationId := ps.ByName("conversationId")
	if !utils.CheckIdentifier(conversationId) {
		ctx.Logger.WithError(err).Error("Error: conversation ID not valid")
		http.Error(w, "Error: conversation ID not valid", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()

	paramLimit := queryParams.Get("limit")
	cursor := queryParams.Get("cursor")

	// Usa i parametri come desiderato
	if paramLimit == "" {
		ctx.Logger.Error("Missing required query parameters")
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}
	if cursor != "" && !utils.CheckIdentifier(cursor) {
		ctx.Logger.Error("cursor not valid")
		http.Error(w, "cursor not valid", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(paramLimit)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid param1: expected an integer")
		http.Error(w, "Invalid param1: expected an integer", http.StatusBadRequest)
		return
	}

	if limit <= 0 {
		ctx.Logger.WithError(err).Error("limit must be positivi and greater then zero")
		http.Error(w, "Invalid limit integer", http.StatusBadRequest)
		return
	}

	response, code, err := rt.db.GetConversation(auth, conversationId, limit, cursor)
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

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("JSON ENCODING")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	err = rt.db.SetRecvMessage(auth, conversationId)
	if err != nil {
		ctx.Logger.WithError(err).Error("update received status")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}
