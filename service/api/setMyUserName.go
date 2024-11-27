package api

import (
	"WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	w.Header().Set("Content-Type", "application/json")
	var requestBody setMyUserNameRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
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

			var responseBody setMyUserNameResponseBody
			responseBody.Username = name
			err = json.NewEncoder(w).Encode(responseBody)

			if err != nil {
				http.Error(w, "an error occurred", http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Error: an error occurred during encoding response")
				return
			} else {
				// no errors!
				return
			}

		}

	}
}

func isAuthorized(id string, auth string) bool {
	return id == auth
}
