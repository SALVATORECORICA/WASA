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
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Convert the string to id
	id_user, err := strconv.Atoi(idOfUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// ExTract the banned id
	photo_id := ps.ByName("photo_id")

	// Convert the string to id
	photo_id, err = strconv.Atoi(photo_id)
	if err != nil {
		http.Error(w, "Error by converting the id of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check of the photo owner´s is the same of the request
	owner, err := rt.db.OwnerPhoto(photo_id)
	if err != nil {
		http.Error(w, "Error by searching the owner of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by searching the owner of the Photo")
		return
	}
	if owner.Id != id_user {
		http.Error(w, "Error: the user is not allowed to do it", http.StatusMethodNotAllowed)
		ctx.Logger.WithError(err).Error("Deleting operation is not allowed")
		return
	}

	// Delete the photo
	err = rt.db.DeletePhoto(photo_id)
	if err != nil {
		http.Error(w, "Error by deleting the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by deleating the Photo from the DB")
		return
	}
	path, err := GetPhotoPath(photo_id)
	if err != nil {
		http.Error(w, "Error by getting the path of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by getting the path of the Photo")
		return
	}
	err = os.Remove(path)
	if err != nil {
		http.Error(w, "Error by deleting the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by deleating the Photo from the Folder")
		return
	}

	// Delete the comments of the photo
	err = rt.db.DeleteCommentPhoto(photo_id)
	if err != nil {
		http.Error(w, "Error by deleting of the comment of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by deleting of the comment of the Photo")
		return
	}
	// Delete the likes of the photo
	err = rt.db.DeleteLikePhoto(photo_id)
	if err != nil {
		http.Error(w, "Error by deleting of the comment of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error:Error by deleting of the comment of the Photo")
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
