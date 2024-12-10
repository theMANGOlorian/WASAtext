package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
)

func (rt *_router) sendImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	if r.ContentLength == 0 || r.Header.Get("Content-Type") != "image/png" {
		ctx.Logger.Error("Error: Content-Type not supported")
		http.Error(w, "Content-Type not valid", http.StatusBadRequest)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: Reading file")
		http.Error(w, "Error: couldn't read the image", http.StatusInternalServerError)
		return
	}

	code, response, err := rt.db.SendImage(auth, conversationId)

	if err != nil {
		if code == 404 {
			ctx.Logger.WithError(err).Error("Error: user/group not found")
			http.Error(w, "Error: user/group not found", http.StatusNotFound)
			return
		}
		if code == 403 {
			ctx.Logger.WithError(err).Error("Error: Operation forbidden")
			http.Error(w, "Error: operation forbidden", http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("An error occurred")
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	const imagesPath = "/tmp/WasaText/images/"
	osFile, err := os.Create(imagesPath + response.Image + ".png")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error: encoding json ")
		http.Error(w, "Error: writing response", http.StatusInternalServerError)
		return
	}

}
