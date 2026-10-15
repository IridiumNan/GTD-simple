package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IridiumNan/GTD-simple/internal/capture/store"
	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/gin-gonic/gin"
)

func handlerCaptureCreate(c *gin.Context) {
	var reqJSON RequestCapture

	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, NewResponse("json struct error", nil))
		return
	}
	var ns *models.Stuff
	if reqJSON.Data.S.ID == "" {
		ns = models.NewStuff(reqJSON.Data.S.Raw)
	} else {
		// use the request structure
		ns = &reqJSON.Data.S
	}

	store.Create(ns)

	c.JSON(http.StatusCreated, NewResponse(
		fmt.Sprintf("stuff has been created, id: %s, content: %s", ns.ID, ns.Raw),
		gin.H{
			"id":   ns.ID,
			"stat": store.Stat(),
		}))
}

// handlerCaptureDelete delete specific stuff by id
// /capture/fetch/:id get id from endpoint
func handlerCaptureDelete(c *gin.Context) {
	idStr := c.Param("id")
	id := strings.TrimSpace(idStr)

	exist := store.Delete(id)
	if !exist {
		c.JSON(http.StatusNotFound, NewResponse(fmt.Sprintf("stuff with id: %s not found, this request is omitted", id), nil))
	}

	c.JSON(http.StatusOK, NewResponse(fmt.Sprintf("stuff has been marked as deleted, id: %s", id),
		gin.H{
			"stat": store.Stat(),
		}))
}

// handlerCaptureFetch fetch a stuff by id
// /capture/delete/:id get id from endpoint
func handlerCaptureFetch(c *gin.Context) {
	idStr := c.Param("id")

	id := strings.TrimSpace(idStr)

	stuf := store.Fetch(id)
	if stuf == nil {
		c.JSON(http.StatusNotFound,
			NewResponse(
				fmt.Sprintf("stuff not found, id: %s", id),
				nil),
		)
		return
	}

	c.JSON(http.StatusOK, NewResponse("stuff found", map[string]any{
		"stuf": stuf,
	}))
}

func handlerCaptureClear(c *gin.Context) {
	store.Clear()

	c.JSON(http.StatusOK, NewResponse("all unused stuffs has been removed",
		gin.H{
			"stat": store.Stat(),
		}))
}
