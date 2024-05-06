package api

func (rt *_router) postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Extracting the id of the user
	idOfUser := extractBearer(r.Header.Get("Authorization"))

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

	// Convert the id from string to int
	idUser, err := strconv.Atoi(idOfUser)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id is valid
	if valid, err := rt.db.SearchUserID(idUser); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Exract and check the photo id
	id_photo := ps.ByName("id_photo")

	// Convert the string to id
	id_photo, err = strconv.Atoi(id_photo)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id of the photo is valid
	if valid, err := rt.db.ExistsPhoto(id_photo); !valid || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Extract the comment from the request body
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("comment: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check of the comment is valid
	valid := validComment(comment.comment)
	if !valid {
		ctx.Logger.Error("comment: error validating comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Save the comment in the db
	err = rt.db.PostComment(id_photo, idUser, comment.comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("Inserting-Comment: error by inserting comment in the db")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
