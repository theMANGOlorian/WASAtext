package api

import (
	"WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	var doLoginRequestBody doLoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&doLoginRequestBody)
	if err != nil {
		// Decoding JSON error
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	} else {
		// check Username field is > 0
		if len(doLoginRequestBody.Username) == 0 {
			ctx.Logger.WithError(err).Error("Error: username length  is zero")
			http.Error(w, "username can not be empty", http.StatusBadRequest)
		} else {
			// there are no errors in the request
			id, err := rt.db.DoLogin(doLoginRequestBody.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Error: an error occurred during database operations")
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				var doLoginResponseBody doLoginResponseBody
				doLoginResponseBody.Identifier = id

				err = json.NewEncoder(w).Encode(doLoginResponseBody)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
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
