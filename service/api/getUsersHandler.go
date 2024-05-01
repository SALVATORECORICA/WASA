package api

import (
	"encoding/json"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) getUsersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extracting the id of the user
	idOfUser := extractBearer(r.Header.Get("Authorization"))
	print("id:", idOfUser)

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
		fmt.Println("ciao")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Extract the query parameter from the URL
	nicknameSearched := r.URL.Query().Get("nickname")

	// Search the user
	users, err := rt.db.SearchUserFromNick(nicknameSearched)
	if err != nil {
		http.Error(w, "Error by searching of user", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Show users who haven't blocked the user searching for profiles
	/*users, err = rt.db.CheckBan(users, idUser)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The check ban has encountered in error")
		return
	}*/
	// return the users
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	return
}
