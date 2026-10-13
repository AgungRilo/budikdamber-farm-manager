// ===== internal/handler/feed_handler.go =====
package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type FeedHandler struct {
	svc *service.FeedService
}

func NewFeedHandler(svc *service.FeedService) *FeedHandler {
	return &FeedHandler{svc: svc}
}

// List godoc
// @Summary   Daftar pakan
// @Tags      Feeds
// @Produce   json
// @Security  BearerAuth
// @Param     page    query  int     false  "Halaman"  default(1)
// @Param     limit   query  int     false  "Jumlah per halaman"  default(10)
// @Param     search  query  string  false  "Cari nama / merek"
// @Success   200  {object}  response.Body{data=[]model.Feed,meta=response.Meta}
// @Router    /feeds [get]
func (h *FeedHandler) List(c *gin.Context) {
	f := model.FeedFilter{
		Search: c.Query("search"),
		Page:   queryInt(c, "page", 1, 1, 100000),
		Limit:  queryInt(c, "limit", 10, 1, 100),
	}
	list, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		handleError(c, err, "pakan")
		return
	}
	response.Paginated(c, "berhasil mengambil daftar pakan", list, response.NewMeta(f.Page, f.Limit, total))
}

// Get godoc
// @Summary   Detail pakan
// @Tags      Feeds
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID pakan"
// @Success   200  {object}  response.Body{data=model.Feed}
// @Failure   404  {object}  response.Body
// @Router    /feeds/{id} [get]
func (h *FeedHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	f, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, "pakan")
		return
	}
	response.OK(c, "berhasil mengambil detail pakan", f)
}

// Create godoc
// @Summary   Buat pakan (owner)
// @Tags      Feeds
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.FeedInput  true  "Data pakan"
// @Success   201   {object}  response.Body{data=model.Feed}
// @Failure   400   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /feeds [post]
func (h *FeedHandler) Create(c *gin.Context) {
	var in model.FeedInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	f, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		handleError(c, err, "pakan")
		return
	}
	response.Created(c, "pakan berhasil dibuat", f)
}

// Update godoc
// @Summary   Ubah pakan (owner)
// @Tags      Feeds
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int              true  "ID pakan"
// @Param     body  body      model.FeedInput  true  "Data pakan"
// @Success   200   {object}  response.Body{data=model.Feed}
// @Failure   404   {object}  response.Body
// @Router    /feeds/{id} [put]
func (h *FeedHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.FeedInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	f, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		handleError(c, err, "pakan")
		return
	}
	response.OK(c, "pakan berhasil diperbarui", f)
}

// Delete godoc
// @Summary   Hapus pakan / soft delete (owner)
// @Tags      Feeds
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID pakan"
// @Success   200  {object}  response.Body
// @Failure   404  {object}  response.Body
// @Router    /feeds/{id} [delete]
func (h *FeedHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err, "pakan")
		return
	}
	response.OK(c, "pakan berhasil dihapus", nil)
}
