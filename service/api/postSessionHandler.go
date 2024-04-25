package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) postSessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

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

	// the structure that take the vaules from thw Json
	type Data struct {
		nickname string `json:"nickname"`
	}

	// We read the nickname
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error by parsing of the JSON", http.StatusBadRequest)
		return
	}

	//We check that the field is not empty
	if data.nickname == "" {
		http.Error(w, "The field username is empty", http.StatusBadRequest)
		return
	}

}

//Check of the id is valid and if it already exists, if not we create a new entry in the db for the table "User" and return a new user_id if yes we return only the id
func checkAndAddUser(string nickname) {
	if !isValidID(nickname) {
		w.WriteHeader(http.StatuBadRequest)
		return
	}
	//
	id, err := putNewUser(nickname)
	if err != nil{
		http.Error(w."Error with the DB", http.StatusBadRequest)
	}

}
