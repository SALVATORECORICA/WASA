package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/Struct"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"time"
)

func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

	//
	var image Struct.Image

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Check of the HTTP method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ctx.Logger.Error("Method is not correct, the method should be POST")
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
	if valid, err := rt.db.SearchUserID(idUser); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// extract the photo from the body request
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert the base64 image in array of byte
	imageBytes, err := base64.StdEncoding.DecodeString(image.Photo_data)
	if err != nil {
		http.Error(w, "failed by decoding the image", http.StatusBadRequest)
		return
	}

	// obtain the nickname to find the folder where we want to save the photo
	nickname, err := SearchNickname(idUser)
	if err != nil {
		http.Error(w, "failed by searching the nickname for the path", http.StatusBadRequest)
		return
	}

	// obtain the path where we want to save the photo
	path, err := os.Executable()
	if err != nil {
		http.Error(w, "failed by obtaining the path", http.StatusBadRequest)
		return
	}

	timestamp := time.Now()
	path = filepath.Dir(path)
	photoDir := "photos"
	path = filepath.Join(path, nickname, photoDir)

	completePath, _, err := rt.db.PostNewPhoto(nickname, path, timestamp)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error by inserting of the photo in the DB")
		http.Error(w, "Error by inserting of the photo in the DB", http.StatusBadRequest)

	}

	// Save the bytes
	file, err := os.Create(completePath)
	if err != nil {
		http.Error(w, "failed to create image file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write the image
	if _, err := file.Write(imageBytes); err != nil {
		http.Error(w, "failed to write image data to file", http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
