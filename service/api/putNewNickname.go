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

	// Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Confirm the identity of the user
	pathId := ps.ByName("id")

	// Convert the string to int
	idUser, err := strconv.Atoi(pathId)
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
	var nick structures.UserNickname
	err = json.NewDecoder(r.Body).Decode(&nick)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search in the db if the Nickname already exists
	exists, err := rt.db.SearchUser(nick.Nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error searching user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if exists != -1 {
		ctx.Logger.WithError(err).Error("update-nickname: the nickname was already assigned")
		http.Error(w, "Operation not permitted", http.StatusForbidden)
	}

	// Update the nickname
	err = rt.db.PutNewNickname(nick.Nickname, idUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error by the updating of the nickname")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
	}

}
