package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/capture/create", handlerCaptureCreate)
	r.DELETE("/capture/delete/:id", handlerCaptureDelete)
	r.GET("/capture/fetch/:id", handlerCaptureFetch)
	r.DELETE("/capture/clear", handlerCaptureClear)

	r.Static("/static", "./front")

	return r
}
