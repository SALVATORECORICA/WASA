package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa-1967862/service/api/reqcontext"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) putLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Extracting the id of the user
	idOfUser := extractBearer(r.Header.Get("Authorization"))

	// If the user is not logged in then respond with a 403 http status
	if idOfUser == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check of the path id correspond to the beaerer
	pathId := ps.ByName("id")
	if pathId != idOfUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Convert the id from string to int
	idUser, err := strconv.Atoi(idOfUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id is valid
	if valid, err := rt.db.ExistsUser(idUser); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Exract and check the photo
	photoId := ps.ByName("photo_id")

	// Convert the string to id
	photoIdInt, err := strconv.Atoi(photoId)
	if err != nil {
		http.Error(w, "Error by converting the id of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Search in the DB of the photo id is valid
	if valid, err := rt.db.ExistsPhoto(photoIdInt); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Prendo il proprietario della foto, controllo se lui mi ha bannato
	ownerOfPhoto, err := rt.db.OwnerPhotoFromIdPhoto(photoIdInt)
	if err != nil {
		http.Error(w, "Error by converting searching the Photo Owner", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	existsBan, err := rt.db.ExistsBan(ownerOfPhoto.Id, idUser)
	if err != nil {
		http.Error(w, "Error by converting searching the  ban from the Photo Owner", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if existsBan {
		http.Error(w, "Operation not permitted, the user was banned ", http.StatusMethodNotAllowed)
		ctx.Logger.WithError(err).Error("The user was banned and cannot put the like")
		return
	}

	// Check of the like already exists
	existsLike, err := rt.db.ExistsLike(idUser, photoIdInt)
	if err != nil {
		http.Error(w, "Error by searching if the like already exists", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by searching if the like already exists")
		return
	}
	if existsLike {
		http.Error(w, "Operation not permitted, the like was already putted", http.StatusMethodNotAllowed)
		ctx.Logger.WithError(err).Error("The like was already putted")
		return
	}
	// Insert the like
	err = rt.db.PutLike(photoIdInt, idUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("inserting-ban: error by the inserting of the ban")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
	return
}
