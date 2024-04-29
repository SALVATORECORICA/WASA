package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extracting the id of the user
	idOfUser := extractBearer(r.Header.Get("Authorization"))
	// If the user is not logged in then respond with a 403 http status
	if idOfUser == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Extract the query parameter from the URL
	nicknameSearched := r.URL.Query().Get("nickname")

	// Search the user
	var users []User
	rows, err := rt.db.Query("SELECT * FROM users WHERE nickname LIKE ?", nicknameSearched+"%")
	if err != nil {
		http.Error(w, "Error by searching of user", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	defer rows.Close()

	//Save the founded users
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Nickname)
		if err != nil {
			http.Error(w, "Error by searching of user", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Database has encountered an error")
			return
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, "Error searching for user", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

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
