package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa-1967862/service/api/reqcontext"
)

func (rt *_router) putFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Check of the Server is ready:
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("The Server is not ready")
		return
	}

	pathId := ps.ByName("id")

	// Convert the string to id
	idUser, err := strconv.Atoi(pathId)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id is valid
	if valid, err := rt.db.ExistsUser(idUser); !valid || err != nil {
		http.Error(w, "Error by searching the User: THe User not exists ", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Exract and check the followed id
	followedId := ps.ByName("followed_id")

	// Convert the string to id
	followedIdInt, err := strconv.Atoi(followedId)
	if err != nil {
		http.Error(w, "Error by converting the id of the User", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	// Search in the DB of the id is valid
	if valid, err := rt.db.ExistsUser(followedIdInt); !valid || err != nil {
		http.Error(w, "Error by searching the following User: TH User not exist", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check if the user wants to follow themselves; it doesn't make sense
	if idUser == followedIdInt {
		http.Error(w, "Error by putting the following", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("The user cannot following themselves")
		return
	}

	// Check if exists a ban
	exists, err := rt.db.ExistsBan(followedIdInt, idUser)
	if err != nil {
		http.Error(w, "Error by searching of the ban", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if exists {
		http.Error(w, "Exists a ban from the user, the user cannot be followed", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check if exists a ban
	exists, err = rt.db.ExistsBan(idUser, followedIdInt)
	if err != nil {
		http.Error(w, "Error by searching of the ban", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}
	if exists {
		http.Error(w, "Exists a ban from the user, the user cannot be followed", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		return
	}

	// Check if the following already exists
	existsFollowing, err := rt.db.ExistsFollowing(idUser, followedIdInt)
	if err != nil {
		http.Error(w, "Error by searching if the following already exists", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error by searching if the following already exists")
		return
	}
	if existsFollowing {
		http.Error(w, "Operation not permitted, the following was already putted", http.StatusMethodNotAllowed)
		ctx.Logger.WithError(err).Error("The following was already putted")
		return
	}

	// Insert the following
	err = rt.db.PutFollowing(idUser, followedIdInt)
	if err != nil {
		ctx.Logger.WithError(err).Error("inserting-ban: error by the inserting of the ban")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
	}
}
