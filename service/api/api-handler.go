package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login
	rt.router.POST("/session", rt.wrap(rt.postSessionHandler))

	//search
	rt.router.GET("/users", rt.wrap(rt.getUsersHandler))

	//user
	// rt.router.GET("/users/:id", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:id", rt.wrap(rt.putNewNickname))

	//ban
	// rt.router.PUT("/users/:id/banned_users/:banned_user_id", rt.wrap(rt.putNewBan))
	// rt.router.DELETE("/users/:id/banned_users/:banned_user_id", rt.wrap(rt.deleteBan))

	//followers
	// rt.router.PUT("/users/:id/followers/:follower_id", rt.wrap(rt.putFollowing))
	// rt.router.DELETE("/users/:id/followers/:follower_id", rt.wrap(rt.deleteFollowing))

	//stream
	// rt.router.GET("/users/:id/home", rt.wrap(rt.getStream))

	//Upload photo
	// rt.router.POST("/users/:id/photos", rt.wrap(rt.postPhoto))

	//photo
	// rt.router.GET("/users/:id/photos/:photo_id", rt.wrap(rt.getPhoto))
	// rt.router.DELETE("/users/:id/photos/:photo_id", rt.wrap(rt.deletePhoto))

	//comments
	// rt.router.POST("users/:id/photos/:photo_id/comments", rt.wrap(rt.postComment))
	// rt.router.DELETE("users/:id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment))

	//likes
	// rt.router.PUT("/users/:id/photos/:photo_id/likes", rt.wrap(rt.putLike))
	// rt.router.DELETE("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.deleteLike))

	//example
	// rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
