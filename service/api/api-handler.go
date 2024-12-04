package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Default Routes (from template)
	rt.router.GET("/", rt.getHelloWorld)
	// rt.router.GET("/context", rt.wrap(rt.getContextReply))

	/*WasaText routes*/
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:userId/photoProfile", rt.wrap(rt.setMyPhoto))
	rt.router.POST("/users/:userId/conversations", rt.wrap(rt.startConversation))
	rt.router.GET("/users/:userId/conversations", rt.wrap(rt.getConversations))
	rt.router.PUT("/conversations/:conversationId/member", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversations/:conversationId", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/conversations/:conversationId/groupName", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversations/:conversationId/groupPhoto", rt.wrap(rt.setGroupPhoto))
	rt.router.POST("/conversations/:conversationId/text", rt.wrap(rt.sendMessage))
	// Special routes (from template)
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
