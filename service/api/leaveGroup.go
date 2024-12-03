package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}
	conversationId := ps.ByName("conversationId")

	if !utils.CheckIdentifier(conversationId) {
		ctx.Logger.WithError(err).Error("Error: conversationId not valid")
		http.Error(w, "Error: conversationId not valid", http.StatusBadRequest)
		return
	}

	code, err := rt.db.LeaveGroup(auth, conversationId)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("group/user not found")
			http.Error(w, "Error: group and/or user don't exists", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("an error occurred")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
