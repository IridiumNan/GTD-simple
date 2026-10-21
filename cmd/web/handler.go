package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/IridiumNan/GTD-simple/internal/capture/store"
	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/gin-gonic/gin"
)

// handlerCaptureCreate creating a new stuff
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
		return
	}

	c.JSON(http.StatusOK, NewResponse(fmt.Sprintf("stuff has been marked as deleted, id: %s", id),
		gin.H{
			"stat": store.Stat(),
		}))
}

func handlerCaptureMultiDel(c *gin.Context) {
	type idList struct {
		List []string `json:"id_list"`
	}
	type multiDelBody struct {
		Data idList `json:"data"`
	}

	var reqJSON multiDelBody
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, NewResponse("json stuct unmarshal error ", nil))
		return
	}

	notFound := make([]string, 0, len(reqJSON.Data.List))
	for _, id := range reqJSON.Data.List {
		exist := store.Delete(id)
		if !exist {
			slog.Warn("soft deleted for not exist stuff", "id", id)
			notFound = append(notFound, id)
		}
	}

	c.JSON(http.StatusOK, NewResponse("all stuffs has been marked as deleted", gin.H{
		"not_found": notFound,
		"stat":      store.Stat(),
	}))
}

func handlerCaptureUpdate(c *gin.Context) {
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
		fmt.Sprintf("stuff has been updated, id: %s, content: %s", ns.ID, ns.Raw),
		gin.H{
			"id":   ns.ID,
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

	c.JSON(http.StatusOK, NewResponse("stuff found", gin.H{
		"stuff": stuf,
	}))
}

// handlerCaptureClear clear all deleted stuffs
// /capture/clear endpoint
func handlerCaptureClear(c *gin.Context) {
	store.Clear()

	c.JSON(http.StatusOK, NewResponse("all unused stuffs has been removed",
		gin.H{
			"stat": store.Stat(),
		}))
}

// handlerCaptureAll get all not deleted stuffs
// /capture/all
func handlerCaptureAll(c *gin.Context) {
	fStr := c.Param("filter")
	fStr = strings.Trim(fStr, "/ ")

	f, ok := store.Filter[fStr]
	if !ok {
		c.JSON(http.StatusBadRequest, NewResponse(fmt.Sprintf("filter %s not found", fStr), nil))
		return
	}

	stuffs := store.All(f)

	if len(stuffs) == 0 {
		c.JSON(http.StatusOK, NewResponse("These is no stuff", gin.H{
			"stuffs": stuffs,
			"stat":   store.Stat(),
		}))
		return
	}

	c.JSON(http.StatusOK, NewResponse("Get all stuffs not deleted success", gin.H{
		"stuffs": stuffs,
		"stat":   store.Stat(),
	}))
}

// handlerCaptureFlush flush current data on ram into disk
func handlerCaptureFlush(c *gin.Context) {
	err := store.Flush()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewResponse(fmt.Sprintf("Flush error: %s", err.Error()), nil))
		return
	}

	c.JSON(http.StatusOK, NewResponse("Flush success", gin.H{
		"stat": store.Stat(),
	}))
}
