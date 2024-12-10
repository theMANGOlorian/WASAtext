package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var requestBody utils.SetMyUserNameRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// Decoding JSON error
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	} else {
		if len(requestBody.Username) == 0 {
			// strings empty are not allowed
			ctx.Logger.WithError(err).Error("Error: username length is zero ")
			http.Error(w, "username as string empty is not allowed", http.StatusBadRequest)
			return
		} else {
			id := ps.ByName("userId")

			if !isAuthorized(id, auth) {
				ctx.Logger.WithError(err).Error("Error: Operation Unauthorized")
				http.Error(w, "Operation Unauthorized", http.StatusUnauthorized)
				return
			}

			name, err := rt.db.SetMyUserName(id, requestBody.Username)
			if err != nil {
				if err.Error() == "error while updating username: UNIQUE constraint failed: users.username" {
					http.Error(w, "Username already exists", http.StatusConflict)
					return
				}
				ctx.Logger.WithError(err).Error("Error: an error occurred during database operation")
				http.Error(w, "an error occurred", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			var responseBody utils.SetMyUserNameResponseBody
			responseBody.Username = name
			err = json.NewEncoder(w).Encode(responseBody)

			if err != nil {
				http.Error(w, "an error occurred", http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Error: an error occurred during encoding response")
				return
			}

		}

	}
}

func isAuthorized(id string, auth string) bool {
	return id == auth
}
