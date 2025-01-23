package api

import (
	"WASAtext/service/api/reqcontext"
	"WASAtext/service/api/utils"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	auth, err := utils.CheckAuthorizationField(authHeader)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "YOU SHALL NOT PASS", http.StatusUnauthorized)
		return
	}

	photoId := ps.ByName("photoId")

	if !utils.CheckPhotoId(photoId) {
		ctx.Logger.WithError(err).Error("Error: photoId not valid")
		http.Error(w, "photo ID not valid", http.StatusBadRequest)
		return
	}

	allowed, err := rt.db.GetPhoto(auth, photoId)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "an error occurred", http.StatusForbidden)
		return
	}

	if !allowed {
		ctx.Logger.Error("operation forbidden")
		http.Error(w, "cannot access to this data", http.StatusForbidden)
		return
	}

	file, err := os.Open("/tmp/WasaText/images/" + photoId + ".png")
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "cannot send file", http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ctx.Logger.Error(err)
		}
	}(file)

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		ctx.Logger.Error(err)
		http.Error(w, "cannot send file", http.StatusInternalServerError)
		return
	}

}
