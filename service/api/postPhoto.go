package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"wasa-1967862/service/api/reqcontext"
)

func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

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
		fmt.Println("user non valido nww")
		return
	}

	// Read the body of the request
	imageBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error by reading the body of the request", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by reading of the body")
		return
	}

	// Check if the image Type is valid
	validType := detectImageType(imageBytes)
	if !validType {
		ctx.Logger.WithError(err).Error("Error: format must be PNG or JPEG")
		http.Error(w, "Error: format must be PNG or JPEG", http.StatusBadRequest)
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
	path = filepath.Join(path, idOfUser, photoDir)
	fmt.Println("il path non completo e ", path)
	_, completePath, err := rt.db.PostNewPhoto(idUser, path, timestamp)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error by inserting of the photo in the DB")
		http.Error(w, "Error by inserting of the photo in the DB", http.StatusBadRequest)
	}
	fmt.Println("il path e", completePath)
	// Save the bytes
	file, err := os.Create(completePath)
	if err != nil {
		http.Error(w, "failed to create image file", http.StatusInternalServerError)
		return
	}
	fmt.Println("il path e", completePath)
	defer file.Close()
	// Write the image
	if _, err := file.Write(imageBytes); err != nil {
		http.Error(w, "failed to write image data to file", http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
