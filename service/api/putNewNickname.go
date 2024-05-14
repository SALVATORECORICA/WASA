package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa-1967862/service/api/reqcontext"
	"wasa-1967862/service/structures"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) putNewNickname(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
		ctx.Logger.Error("The user is not logged")
		return
	}

	// Confirm the identity of the user
	pathId := ps.ByName("id")
	if pathId != idOfUser {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("The user is not allowed to change the nickname of other users")
		return
	}

	// Convert the string to int
	idUser, err := strconv.Atoi(idOfUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Search in the DB of the id is valid
	if valid, err := rt.db.ExistsUser(idUser); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("The bearer and the user id not exist in the db")
		return
	}

	// Get the new nickname from the body
	var nick structures.Data
	err = json.NewDecoder(r.Body).Decode(&nick)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the nickname
	err = rt.db.PutNewNickname(nick.Nickname, idUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error by the updating of the nickname")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
	}

}
