package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/Struct"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
// We have 3 input parameters, the first is the reply of the HTTP Request, the second one is the URL and Body request, the third one is the parameters of the URL Path
func (rt *_router) postSessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Set a reply as JSON
	w.Header().Set("Content-Type", "application/json")

	//Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	// Check of the HTTP method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ctx.Logger.Error("Method is not correct, the method should be POST")
		return
	}

	// the structure that take the values from the Json

	var data Struct.Data
	var Req Struct.DataId

	// We read the nickname
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Error by parsing of the JSON", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by parsing the JSON")
		return
	}

	//Check of the nickname is valid
	if !isValidID(data.Nickname) {
		http.Error(w, "The username is not valid", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("The username is not valid")
		return
	}

	// search the user in the db
	id, err2 := rt.db.SearchUser(data.Nickname)
	if err2 != nil {
		http.Error(w, "Error by DB", http.StatusBadRequest)
		ctx.Logger.WithError(err2).Error("Error by DB")
		return
	}
	// user not in DB --> we create a new user
	if id == -1 {
		id, err := rt.db.PutNewUser(data.Nickname)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error by creating new user")
			http.Error(w, "Error by creating new user", http.StatusBadRequest)
			//
			err = createFolders(data.Nickname)
			if err != nil {
				ctx.Logger.WithError(err).Error("Error by creating new user")
				http.Error(w, "Error by creating the folder of the user", http.StatusBadRequest)
				return
			}
			return
		}

		Req.Id = id
		_ = json.NewEncoder(w).Encode(Req)
		return

	}
	Req.Id = int(id)
	_ = json.NewEncoder(w).Encode(Req)
	return

}
