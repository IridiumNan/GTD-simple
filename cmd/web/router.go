package web

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IridiumNan/GTD-simple/front"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func buildfrontFS() fs.FS {
	if devDir := viper.GetString("service.front.dir"); devDir != "" {
		// for develop mode, read the dir directly, no need to rebuild this cli
		slog.Info("using the front dir directly (Only do this for development)", "path", devDir)
		return os.DirFS(devDir)
	}

	slog.Info("using the embeded frontend source")
	return front.FS
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Router() *gin.Engine {
	frontFS := buildfrontFS()
	r := gin.Default()
	r.Use(timeoutMiddleware(2 * time.Second))
	// register backend endpoints
	r.POST("/capture/create", handlerCaptureCreate)

	r.DELETE("/capture/delete/:id", handlerCaptureDelete)
	r.DELETE("/capture/multiDel", handlerCaptureMultiDel)

	r.POST("/capture/update", handlerCaptureUpdate)

	r.GET("/capture/fetch/:id", handlerCaptureFetch)

	r.DELETE("/capture/clear", handlerCaptureClear)

	r.GET("/capture/all/*filter", handlerCaptureAll)

	r.POST("/capture/flush", handlerCaptureFlush)

	frontServe := gin.WrapH(http.FileServerFS(frontFS))

	// static files
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		frontServe(c)
	})

	return r
}

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
