package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
	var doLoginRequestBody utils.DoLoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&doLoginRequestBody)
	if err != nil {
		// Decoding JSON error
		ctx.Logger.WithError(err).Error("Error: Decoding JSON ")
		http.Error(w, "Cannot parse RequestBody", http.StatusBadRequest)
		return
	}

	if !utils.CheckName(doLoginRequestBody.Username) {
		ctx.Logger.WithError(err).Error("Error: username not valid")
		http.Error(w, "username not valid, at least 3 characters and less then 26", http.StatusBadRequest)
		return
	}

	id, photoCode , err := rt.db.DoLogin(doLoginRequestBody.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error: an error occurred during database operations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var doLoginResponseBody utils.DoLoginResponseBody
	doLoginResponseBody.Identifier = id
	doLoginResponseBody.Username = doLoginRequestBody.Username
	doLoginResponseBody.PhotoCode = photoCode

	err = json.NewEncoder(w).Encode(doLoginResponseBody)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error: an error occurred during encoding response")
		return
	}

}
