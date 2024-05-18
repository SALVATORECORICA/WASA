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
func (rt *_router) putNewBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Check of the HTTP method is PUT
	if r.Method != http.MethodPut {
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
	if valid, err := rt.db.ExistsUser(idUser); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Exract and check the banned id
	banned_id := ps.ByName("banned_user_id")

	// Convert the string to id
	banned_idInt, err := strconv.Atoi(banned_id)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id is valid
	if valid, err := rt.db.ExistsUser(banned_idInt); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// Insert the Ban
	err = rt.db.PutNewBan(idUser, banned_idInt)
	if err != nil {
		ctx.Logger.WithError(err).Error("inserting-ban: error by the inserting of the ban")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
	}
}
