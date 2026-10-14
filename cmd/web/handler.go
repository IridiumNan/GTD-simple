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

	ns := models.NewStuff(reqJSON.Data.S.Raw)
	store.Create(ns)

	c.JSON(http.StatusCreated, NewResponse(
		fmt.Sprintf("stuff has been created, id: %s, content: %s", ns.ID, ns.Raw),
		map[string]string{
			"id": ns.ID,
		}))
}

func handlerCaptureDelete(c *gin.Context) {
	idStr := c.Param("id")
	id := strings.TrimSpace(idStr)

	exist := store.Delete(id)
	if !exist {
		c.JSON(http.StatusNotFound, NewResponse(fmt.Sprintf("stuff with id: %s not found, this request is omitted", id), nil))
	}

	c.JSON(http.StatusOK, NewResponse(fmt.Sprintf("stuff with id %s has been marked as deleted", id), nil))
}

// func handlerCaptureFetch(c *gin.Context) {
// 	idStr := c.Param("id")
//
// 	// TODO:
// }
