package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Check of the HTTP method is GET
	if r.Method != http.MethodGet {
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

	// Convert the string to id
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

	// Extract id of the requester
	idProfileSearched := ps.ByName("id")

	// We verify that the profile searched by the user has not banned him

	exists, err := rt.db.ExistsBan(idProfileSearched, idUser)
	if err != nil {
		http.Error(w, "Error by searching the ban", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: By Searching the Ban")
		return
	}
	if exists {
		http.Error(w, "Error: the user is banned", http.StatusMethodNotAllowed)
		ctx.Logger.WithError(err).Error("It's not permitted to view this profile")
		return
	}

	// We take now the information for the profile

	// Obtain the Nickname
	nickname, err := rt.db.GetNickname(idProfileSearched)
	if err != nil {
		http.Error(w, "Error by getting the Nickname", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: By getting the Nickname")
		return
	}

	// Obtain the Followers
	followers, nFollowers, err := rt.db.GetFollower(idProfileSearched)
	if err != nil {
		http.Error(w, "Error by getting the followers", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: By getting the followers")
		return
	}

	// Obtain the followed
	followed, nFollowed, err := rt.db.GetFollowed(idProfileSearched)
	if err != nil {
		http.Error(w, "Error by getting the followed", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: By getting the followed")
		return
	}

	// Obtain the photos of the profile searched (sorted chronologically)

	photos, err := rt.db.GetPhotosProfileSorted(idProfileSearched)
	if err != nil {
		http.Error(w, "Error by getting the photos", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error: By getting the photos")
		return
	}

	// User Profile
	var userProfile User_Profile

	userProfile.Id = idProfileSearched
	userProfile.Nickname = nickname
	userProfile.Followers = followers
	userProfile.NFollowers = nFollowers
	userProfile.Followed = followed
	userProfile.NFollowed = nFollowed
	userProfile.Photos = photos


	userProfileJSON, err := json.Marshal(photo)
	if err != nil {
		http.Error(w, "Error by creating the JSON, http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error by creating the JSON")
		return
	}

	// Set the Header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON in the reply
	_, err = w.Write(userProfileJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error by writing the JSON")
	}
}

