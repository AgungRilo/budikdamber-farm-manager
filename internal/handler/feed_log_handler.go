package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type FeedLogHandler struct {
	svc *service.FeedLogService
}

func NewFeedLogHandler(svc *service.FeedLogService) *FeedLogHandler {
	return &FeedLogHandler{svc: svc}
}

// cycleLogError memetakan error log (pakan & kematian) ke status HTTP.
func cycleLogError(c *gin.Context, err error, entity string) {
	var aliveErr *repository.AliveExceededError
	if errors.As(err, &aliveErr) {
		c.JSON(http.StatusConflict, response.Body{Success: false, Message: aliveErr.Error(), Data: aliveErr})
		return
	}

	switch {
	case errors.Is(err, repository.ErrCycleNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, repository.ErrCycleClosed):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, repository.ErrFeedNotFound),
		errors.Is(err, repository.ErrLogBeforeStart),
		errors.Is(err, service.ErrLogDateInFuture):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		handleError(c, err, entity)
	}
}

// parseLogFilter membaca ?date_from=&date_to=&page=&limit= ; false kalau format tanggal salah.
func parseLogFilter(c *gin.Context, cycleID int64) (model.LogFilter, bool) {
	f := model.LogFilter{
		CycleID:  cycleID,
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     queryInt(c, "page", 1, 1, 100000),
		Limit:    queryInt(c, "limit", 20, 1, 100),
	}
	for _, d := range []string{f.DateFrom, f.DateTo} {
		if d == "" {
			continue
		}
		if _, err := time.Parse("2006-01-02", d); err != nil {
			response.Error(c, http.StatusBadRequest, "format tanggal filter harus YYYY-MM-DD")
			return f, false
		}
	}
	return f, true
}

// List godoc
// @Summary   Daftar log pakan per siklus
// @Tags      Feed Logs
// @Produce   json
// @Security  BearerAuth
// @Param     id         path   int     true   "ID siklus"
// @Param     date_from  query  string  false  "Dari tanggal (YYYY-MM-DD)"
// @Param     date_to    query  string  false  "Sampai tanggal (YYYY-MM-DD)"
// @Param     page       query  int     false  "Halaman"  default(1)
// @Param     limit      query  int     false  "Jumlah per halaman"  default(20)
// @Success   200  {object}  response.Body{data=[]model.FeedLog,meta=response.Meta}
// @Failure   404  {object}  response.Body
// @Router    /cycles/{id}/feed-logs [get]
func (h *FeedLogHandler) List(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	f, ok := parseLogFilter(c, cycleID)
	if !ok {
		return
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		cycleLogError(c, err, "log pakan")
		return
	}
	response.Paginated(c, "berhasil mengambil log pakan", list, response.NewMeta(f.Page, f.Limit, total))
}

// Create godoc
// @Summary   Catat pemberian pakan (owner & worker)
// @Tags      Feed Logs
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                 true  "ID siklus"
// @Param     body  body      model.FeedLogInput  true  "Data log pakan"
// @Success   201   {object}  response.Body{data=model.FeedLog}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /cycles/{id}/feed-logs [post]
func (h *FeedLogHandler) Create(c *gin.Context) {
	cycleID, ok := parseID(c)
	if !ok {
		return
	}
	var in model.FeedLogInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	l, err := h.svc.Create(c.Request.Context(), cycleID, in, c.GetInt64(middleware.CtxUserID))
	if err != nil {
		cycleLogError(c, err, "log pakan")
		return
	}
	response.Created(c, "log pakan berhasil dicatat", l)
}

// Delete godoc
// @Summary   Hapus log pakan (owner)
// @Tags      Feed Logs
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID log pakan"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Failure   409  {object}  response.Body
// @Router    /feed-logs/{id} [delete]
func (h *FeedLogHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		cycleLogError(c, err, "log pakan")
		return
	}
	response.OK(c, "log pakan berhasil dihapus", nil)
}
