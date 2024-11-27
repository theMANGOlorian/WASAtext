package api

import (
	"WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Authorization
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

	// checking permission
	id := ps.ByName("userId")
	if !isAuthorized(id, auth) {
		ctx.Logger.Error("Error: User not Authorized")
		http.Error(w, "Operation Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.ContentLength == 0 || r.Header.Get("Content-Type") != "image/png" {
		ctx.Logger.Error("Error: Content-Type not supported")
		http.Error(w, "Content-Type not valid", http.StatusBadRequest)
		return
	}

	// Read body request
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Reading file")
		http.Error(w, "Error: couldn't read the image", http.StatusInternalServerError)
		return
	}

	imgCode, err := rt.db.SetMyPhoto(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: generation/saving image code in database")
		http.Error(w, "Error: saving image", http.StatusInternalServerError)
		return
	}

	const imagesPath = "/tmp/WasaText/images/"
	osFile, err := os.Create(imagesPath + imgCode + ".png")
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: saving image")
		http.Error(w, "Error: saving image", http.StatusInternalServerError)
		return
	}
	defer osFile.Close()

	_, err = osFile.Write(buf)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: writing data in png file")
		http.Error(w, "Error: saving image", http.StatusInternalServerError)
		return
	}

	// Generating response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var responseBody setMyPhotoResponseBody
	responseBody.ImageCode = imgCode
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: encoding json ")
		http.Error(w, "Error: writing response", http.StatusInternalServerError)
		return
	}
	return

}
