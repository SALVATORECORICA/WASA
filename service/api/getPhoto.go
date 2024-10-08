package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
	"wasa-1967862/service/api/reqcontext"
	"wasa-1967862/service/structures"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Convert the string to id
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

	// Extract Photo_id
	id := ps.ByName("photo_id")

	// Convert the string to id
	photoId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search the owner of the photo
	owner, err := rt.db.OwnerPhotoFromIdPhoto(photoId)
	if err != nil {
		http.Error(w, "Error by searching the User Owner of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Check of the ban exists
	exists, err := rt.db.ExistsBan(owner.Id, idUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if exists {
		http.Error(w, "The user who search to obtain a photo was banned", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("The user who search to obtain a photo was banned")
		return
	}

	// Now we take the infos of the photo

	// obtain the likes
	usersLikes, nLikes, err := rt.db.GetLikes(photoId)
	if err != nil {
		http.Error(w, "Error by obtaining the likes of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by obtaining the likes of the Photo")
		return
	}
	// obtain the comments
	comments, err := rt.db.CommentsPhoto(photoId, idUser)
	if err != nil {
		http.Error(w, "Error by obtaining the comments of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by obtaining the comments of the Photo")
		return
	}

	// Obtain the date
	date, err := rt.db.GetPhotoDate(photoId)
	if err != nil {
		http.Error(w, "Error by obtaining the date of the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by obtaining the date of the Photo")
		return
	}
	// Obtain the image

	// Obtain the path where we can find the photo
	path, err := rt.db.GetPhotoPath(photoId)
	if err != nil {
		http.Error(w, "Error by obtaining the date of the Path", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by obtaining the date of the Path")
		return
	}

	// Take the photo

	photoData, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Error by uploading the Photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by Uploading the Photo")
		return
	}

	// Check if the user liked the photo
	existsLike, err := rt.db.ExistsLike(idUser, photoId)
	if err != nil {
		http.Error(w, "Error by searching if the like of the photo exists", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: by Searching of the Likes exists")
		return
	}
	// Now we are ready to send the
	var photo structures.Photo
	photo.PhotoId = photoId
	photo.Owner = owner
	photo.Date = date
	photo.Likes = usersLikes
	photo.NLikes = nLikes
	photo.PhotoData = photoData
	photo.Comments = comments
	photo.Liked = existsLike

	// We send the Photo
	photoJSON, err := json.Marshal(photo)
	if err != nil {
		http.Error(w, "Error by creating the JSON", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error by creating the JSON")
		return
	}

	// Set the Header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON in the reply
	_, err = w.Write(photoJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error by writing the JSON")
	}
}
