package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasa-1967862/service/api/reqcontext"
	"wasa-1967862/service/structures"
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

	var userNickname structures.UserNickname
	var userId structures.UserId

	// We read the nickname
	err := json.NewDecoder(r.Body).Decode(&userNickname)

	if err != nil {
		http.Error(w, "Error by parsing of the JSON", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by parsing the JSON")
		return
	}

	//Check of the nickname is valid
	if !isValidID(userNickname.Nickname) {
		http.Error(w, "The username is not valid", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("The username is not valid")
		return
	}

	// search the user in the db
	fmt.Println("entro in search")
	id, err2 := rt.db.SearchUser(userNickname.Nickname)
	fmt.Println("esco")
	if err2 != nil {
		http.Error(w, "Error by DB", http.StatusBadRequest)
		ctx.Logger.WithError(err2).Error("Error by DB")
		return
	}
	// user not in DB --> we create a new user
	if id == -1 {
		fmt.Println("utente non trovato")
		// user not in DB --> we create a new user
		id, err := rt.db.PutNewUser(userNickname.Nickname)
		fmt.Println("id dopo averlo inserito nel db", id)
		if err != nil {
			fmt.Println("entro nell errore")
			ctx.Logger.WithError(err).Error("Error by creating new user")
			http.Error(w, "Error by creating new user", http.StatusBadRequest)
		}
		err = createFolders(id)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error by creating new user")
			http.Error(w, "Error by creating the folder of the user", http.StatusBadRequest)
			return
		}
		userId.Id = id
		_ = json.NewEncoder(w).Encode(userId)
		return

	}
	userId.Id = id
	fmt.Println(userId.Id)
	_ = json.NewEncoder(w).Encode(userId)
	return

}
