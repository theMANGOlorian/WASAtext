package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")

	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	// checking permission
	id := ps.ByName("userId")
	if !isAuthorized(id, auth) {
		ctx.Logger.Error("Error: User not Authorized")
		http.Error(w, "Operation Forbidden", http.StatusForbidden)
		return
	}

	var responseBody utils.GetConversationsResponseBody
	err = rt.db.GetConversations(id, &responseBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error creating responseBody")
		http.Error(w, "Error creating responseBody", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Encoding JSON ")
		http.Error(w, "Cannot encode JSON", http.StatusInternalServerError)
	}

	return

}
