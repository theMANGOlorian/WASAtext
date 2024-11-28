package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Deafult Routes (from template)
	rt.router.GET("/", rt.getHelloWorld)
	// rt.router.GET("/context", rt.wrap(rt.getContextReply))

	/*WasaText routes*/
	rt.router.POST("/session", rt.wrap(rt.doLogin))                               // doLogin
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.setMyUserName))           // setMyUserName
	rt.router.PUT("/users/:userId/photoProfile", rt.wrap(rt.setMyPhoto))          // setMyPhoto
	rt.router.POST("/users/:userId/conversations", rt.wrap(rt.startConversation)) // startConversation
	// Special routes (from template)
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
