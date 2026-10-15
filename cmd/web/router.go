package web

import (
	"github.com/gin-gonic/gin"
)

// BackRouter return the router for backend endpoints
func BackRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/capture/create", handlerCaptureCreate)

	r.DELETE("/capture/delete/:id", handlerCaptureDelete)

	r.POST("/capture/update", handlerCaptureUpdate)

	r.GET("/capture/fetch/:id", handlerCaptureFetch)

	r.DELETE("/capture/clear", handlerCaptureClear)

	r.GET("/capture/all/*filter", handlerCaptureAll)

	r.POST("/capture/flush", handlerCaptureFlush)

	return r
}
