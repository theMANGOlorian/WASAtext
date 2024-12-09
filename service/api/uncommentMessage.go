package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	code, err := rt.db.UncommentMessage(auth, messageId)
	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("not found")
			http.Error(w, "not found", http.StatusInternalServerError)
			return
		}
		ctx.Logger.WithError(err).Error("an error occurred")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
