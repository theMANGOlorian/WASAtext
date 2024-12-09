package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) readMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	code, err := rt.db.ReadMessage(auth, messageId)
	if code == 500 {
		ctx.Logger.Error(err)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}
	if code == 404 {
		ctx.Logger.Error("not found")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}
