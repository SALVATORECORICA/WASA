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
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Extracting the id of the user
	idUser := extractBearer(r.Header.Get("Authorization"))

	// If the user is not logged in then respond with a 403 http status
	if idUser == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// Check of the path id correspond to the beaerer
	pathId := ps.ByName("id")
	if pathId != idUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Convert the string to id
	idUserInt, err := strconv.Atoi(idUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Extract the comment id
	commentId := ps.ByName("comment_id")

	// Convert the comment id to int
	commentIdInt, err := strconv.Atoi(commentId)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check of the comment id is present in the db
	exists, err := rt.db.ExistsComment(commentIdInt)
	if err != nil {
		http.Error(w, "Error by searching of the comment", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if !exists {
		http.Error(w, "The comment not exists also can not deleted", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Extracting the id of the Photo
	photoId := ps.ByName("photo_id")

	// Convert the comment id to int
	photoIdInt, err := strconv.Atoi(photoId)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check of the user is the owner of the comment
	commentOwner, err := rt.db.OwnerComment(commentIdInt, idUserInt)
	if err != nil {
		http.Error(w, "Error by deleting of the comment, the User is not allowed to do this action", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Check of the user is the owner of the photo
	photoOwner, err := rt.db.OwnerPhotoFromIdPhoto(photoIdInt)
	if err != nil {
		http.Error(w, "Error by deleting of the comment, the User is not allowed to do this action", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	if photoOwner.Id == idUserInt || commentOwner {
		// Delete comment
		err = rt.db.DeleteComment(commentIdInt)
		if err != nil {
			http.Error(w, "Error by deleting of the comment", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Database has encountered an error")
			return
		}
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// To delete the comment the user must be the owner of the comment or the owner of the photo
	http.Error(w, "Operation not permitted", http.StatusForbidden)
	ctx.Logger.WithError(err).Error("Database has encountered an error")

}
