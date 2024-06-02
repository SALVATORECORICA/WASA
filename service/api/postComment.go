package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa-1967862/service/api/reqcontext"
	"wasa-1967862/service/structures"
)

func (rt *_router) postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

	// Check of the Server is ready:
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

	// Exract and check the photo id
	idPhoto := ps.ByName("photo_id")

	// Convert the string to int
	idPhotoInt, err := strconv.Atoi(idPhoto)
	if err != nil {
		http.Error(w, "Error by converting the id of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id of the photo is valid
	if valid, err := rt.db.ExistsPhoto(idPhotoInt); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// we search the id of the owner
	user, err := rt.db.OwnerPhotoFromIdPhoto(idPhotoInt)
	if err != nil {
		http.Error(w, "Error by searching the User from the ID-photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by searching the User from the ID-photo")
	}

	// We check if exists a ban beetwen the photo owner and the comment owner

	existsBan, err := rt.db.ExistsBan(user.Id, idUser)
	if err != nil {
		http.Error(w, "Error by checking the ban", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by by checking the ban")
	}
	if existsBan {
		http.Error(w, "Operation not permitted", http.StatusForbidden)
		return
	}

	// Extract the comment from the request body
	var comment structures.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("comment: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check of the comment is valid
	valid := validComment(comment.Comment)
	if !valid {
		ctx.Logger.Error("comment: error validating comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Save the comment in the db
	err = rt.db.PostComment(idPhotoInt, idUser, comment.Comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("Inserting-Comment: error by inserting comment in the db")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
