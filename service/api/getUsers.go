package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUsersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	usernames, err := rt.db.GetUsersList(auth)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "An error occured", http.StatusInternalServerError)
		return
	}

	var response utils.UsersListResponseBody
	response.Users = *usernames
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "An error occured", http.StatusInternalServerError)
		return
	}
	
}
