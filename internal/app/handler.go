package app

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

const sessionLimit = 16

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests.")
}

func HandleServer(server *gin.Engine) {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 10,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	idList := make([]string, 8)
	sessionMap := make(map[string]*session)

	// Teapot
	server.GET("/", mw, func(c *gin.Context) {
		c.String(http.StatusTeapot, "I'm not a teapot.")
	})

	// New trickle session
	server.POST("/trickle/new-session", mw, func(c *gin.Context) {
		var session session
		for {
			session = newSession()
			if _, ok := sessionMap[session.id]; !ok {
				break
			}
		}
		sessionMap[session.id] = &session
		idList = append(idList, session.id)

		if len(idList) > sessionLimit {
			delete(sessionMap, idList[0])
			idList = idList[1:]
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"id": session.id,
		})
	})

	// Get inviter in trickle session
	server.POST("/trickle/get-inviter", mw, func(c *gin.Context) {
		var req sessionLookupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok": false,
			})
			return
		}

		session, ok := sessionMap[req.ID]
		if !ok || len(session.inviter) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"ok": false,
			})
			return
		}

		info := session.inviter[0]
		session.inviter = session.inviter[1:]
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"inviter": info,
		})
	})

	// Get invitee in trickle session
	server.POST("/trickle/get-invitee", mw, func(c *gin.Context) {
		var req sessionLookupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok": false,
			})
			return
		}

		session, ok := sessionMap[req.ID]
		if !ok || len(session.invitee) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"ok": false,
			})
			return
		}

		info := session.invitee[0]
		session.invitee = session.invitee[1:]

		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"invitee": info,
		})
	})

	// Set description in trickle session
	server.POST("/trickle/set", mw, func(c *gin.Context) {
		var req sessionSetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok": false,
			})
			return
		}

		session, ok := sessionMap[req.ID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"ok": false,
			})
			return
		}

		if req.Inviter != nil {
			session.inviter = append(session.inviter, *req.Inviter)
		}
		if req.Invitee != nil {
			session.invitee = append(session.invitee, *req.Invitee)
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})

	// Delete trickle session
	server.POST("/trickle/delete-session", mw, func(c *gin.Context) {
		var req sessionSetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok": false,
			})
			return
		}

		delete(sessionMap, req.ID)
		for i, id := range idList {
			if id == req.ID {
				idList[i] = idList[len(idList)-1]
				idList = idList[:len(idList)-1]
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})
}
