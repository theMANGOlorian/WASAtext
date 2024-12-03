package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	conversationId := ps.ByName("conversationId")

	var requestBody utils.SetGroupNameRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "json Deconding", http.StatusBadRequest)
		return
	}

	if !utils.CheckIdentifier(conversationId) {
		ctx.Logger.WithError(err).Error("Error: conversationId not valid")
		http.Error(w, "Error: conversationId not valid", http.StatusBadRequest)
		return
	}

	if !utils.CheckName(requestBody.Name) {
		ctx.Logger.WithError(err).Error("Error: requestBody.Name not valid")
		http.Error(w, "Error: Name not valid", http.StatusBadRequest)
		return
	}

	code, err := rt.db.SetGroupName(auth, conversationId, requestBody.Name)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("group/user not found")
			http.Error(w, "group/user not found", http.StatusNotFound)
			return
		}
		if code == 403 {
			ctx.Logger.WithError(err).Error("user cannot modify the resource")
			http.Error(w, "user cannot modify the group of the name", http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var responseBody utils.SetGroupNameResponseBody
	responseBody.Name = requestBody.Name
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		ctx.Logger.WithError(err)
		http.Error(w, "json Deconding", http.StatusInternalServerError)
		return
	}

	return

}
