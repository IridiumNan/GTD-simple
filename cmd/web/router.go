package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// BackRouter return the router for backend endpoints
func BackRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:8000"},

		AllowOrigins: []string{"*"},

		AllowMethods: []string{
			"GET", "POST", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/capture/create", handlerCaptureCreate)

	r.DELETE("/capture/delete/:id", handlerCaptureDelete)
	r.DELETE("/capture/multiDel", handlerCaptureMultiDel)

	r.POST("/capture/update", handlerCaptureUpdate)

	r.GET("/capture/fetch/:id", handlerCaptureFetch)

	r.DELETE("/capture/clear", handlerCaptureClear)

	r.GET("/capture/all/*filter", handlerCaptureAll)

	r.POST("/capture/flush", handlerCaptureFlush)

	return r
}
