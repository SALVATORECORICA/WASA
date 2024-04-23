package api

import (
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) postSessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

	nickname, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	//Check of the Server is ready:
	if err := rt.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check of the HTTP method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nickname, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//Check of the id is valid and if it already exists, if not we create a new entry in the db for the table "User" and "Session" if yes we create a new entry only for the table "session"

func checkAndAddUser(string nickname) {
	if !isValidID(nickname) {
		w.WriteHeader(http.StatuBadRequest)
		return
	}
	_, err := putNewUser(nickname)

}
