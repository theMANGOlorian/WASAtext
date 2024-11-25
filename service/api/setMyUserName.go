package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
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

	token := strings.Split(authHeader, " ")

	claims, err := utils.GetTokenInfo(token[1])
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: token not valid ")
		http.Error(w, "Token not valid", http.StatusUnauthorized)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		ctx.Logger.WithError(err).Error("Error: claims sub not found ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	var requestBody setMyUserNameRequestBody
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
			parts := strings.Split(path, "/")
			if len(parts) != 4 {
				// request path not valid
				ctx.Logger.WithError(err).Error("Error: request path not valid")
				http.Error(w, "request path not valid", http.StatusBadRequest)
				return
			} else {
				id := parts[2]

				if !isAuthorized(id, sub) {
					ctx.Logger.WithError(err).Error("Error: Operation Unauthorized")
					http.Error(w, "Operation Unauthorized", http.StatusUnauthorized)
					return
				}

				name, err := rt.db.SetMyUserName(id, requestBody.Username)
				if err != nil {
					ctx.Logger.WithError(err).Error("Error: an error occurred during database operation")
					http.Error(w, "an error occurred", http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)

				var responseBody setMyUserNameRequestBody
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
}

func isAuthorized(id string, token string) bool {
	return id == token
}
