package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
)

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return 0, false
	}
	return id, true
}

func queryInt(c *gin.Context, key string, def, min, max int) int {
	v, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func queryBool(c *gin.Context, key string) *bool {
	v, err := strconv.ParseBool(c.Query(key))
	if err != nil {
		return nil
	}
	return &v
}

// handleError memetakan error repository ke status HTTP.
func handleError(c *gin.Context, err error, entity string) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Error(c, http.StatusNotFound, entity+" tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		response.Error(c, http.StatusConflict, "nama "+entity+" sudah dipakai")
	case errors.Is(err, repository.ErrInvalid):
		response.Error(c, http.StatusBadRequest, "data "+entity+" tidak valid")
	default:
		log.Printf("%s error: %v", entity, err)
		response.Error(c, http.StatusInternalServerError, "terjadi kesalahan server")
	}
}
