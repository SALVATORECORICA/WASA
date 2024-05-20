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
func (rt *_router) deleteFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Extracting the id of the user
	idOfUser := extractBearer(r.Header.Get("Authorization"))
	print("id:", idOfUser)

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
	idUser, err := strconv.Atoi(idOfUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// ExTract the followed id
	followedId := ps.ByName("followed_id")

	// Convert the string to id
	followedIdInt, err := strconv.Atoi(followedId)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	exists, err := rt.db.ExistsFollowing(idUser, followedIdInt)
	if err != nil {
		http.Error(w, "Error by searching of the following", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if !exists {
		http.Error(w, "The following not exists also can not deleted", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	} else {
		err = rt.db.DeleteFollowing(idUser, followedIdInt)
		if err != nil {
			http.Error(w, "Error by deleting of the ban", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Database has encountered an error")
			return
		}
	}

}
