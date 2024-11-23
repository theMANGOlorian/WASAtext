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

	// dologin //
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Special routes (from template)
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
