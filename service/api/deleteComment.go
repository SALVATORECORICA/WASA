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

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Check of the HTTP method is DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ctx.Logger.Error("Method is not correct, the method should be DELETE")
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

	// Convert the string to int
	idUserInt, err := strconv.Atoi(idUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Extract the comment id
	comment_id := ps.ByName("comment_id")

	// Convert the comment id to int
	comment_idInt, err := strconv.Atoi(comment_id)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check of the comment id is present in the db
	exists, err := rt.db.ExistsComment(comment_idInt)
	if err != nil {
		http.Error(w, "Error by searching of the ban", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if !exists {
		http.Error(w, "The comment not exists also can not deleted", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Extracting the id of the Photo
	photo_id := ps.ByName("photo_id")

	// Convert the comment id to int
	photo_idInt, err = strconv.Atoi(photo_id)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check of the comment id is present in the db
	exists, err = rt.db.ExistsPhoto(photo_idInt)
	if err != nil {
		http.Error(w, "Error by searching of the photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if !exists {
		http.Error(w, "The photo not exists also can the comment not be deleted", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Check of the user is the owner of the comment
	allowedComment, err := rt.db.OwnerComment(idUserInt, comment_idInt)
	if err != nil {
		http.Error(w, "Error by deleting of the comment, the User is not allowed to do this action", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Check of the user is the owner of the photo
	allowedOwner, err := rt.db.OwnerPhoto(idUserInt, comment_idInt)
	if err != nil {
		http.Error(w, "Error by deleting of the comment, the User is not allowed to do this action", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if !allowedComment && !allowedOwner {
		http.Error(w, "The photo not exists also can the comment not be deleted", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Delete comment
	err = rt.db.DeleteComment(comment_idInt)
	if err != nil {
		http.Error(w, "Error by deleting of the comment", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

}
