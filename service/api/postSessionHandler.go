package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) postSessionHandler(w http.ResponseWriter, r *http.Request) {

	// Set a reply as JSON
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

	// the structure that take the values from the Json
	type Data struct {
		nickname string `json:"nickname"`
	}

	var data Data

	// We read the nickname
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error by parsing of the JSON", http.StatusBadRequest)
		return
	}

	//Check of the id is valid
	if !isValidID(nickname) {
		http.Error(w, "The username is not valid", http.StatusBadRequest)
		return
	}

	// search the user in the db
	id, err := rt.db.searchUser(nickname)
	if err != nil {
		w.WriteHeader(http.StatuBadRequest)
		return
	}
	if id == 0 {
		&data, err := putNewUser(nickname)
		if err != nil {
			w.WriteHeader(http.StatuBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
		return
	}
	data.nickname= id
	_ = json.NewEncoder(w).Encode(data)
	return
}
