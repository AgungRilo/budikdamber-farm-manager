package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type PondHandler struct {
	svc *service.PondService
}

func NewPondHandler(svc *service.PondService) *PondHandler {
	return &PondHandler{svc: svc}
}

func (h *PondHandler) List(c *gin.Context) {
	f := model.PondFilter{
		Search:   c.Query("search"),
		Type:     c.Query("type"),
		IsActive: queryBool(c, "is_active"),
		Page:     queryInt(c, "page", 1, 1, 100000),
		Limit:    queryInt(c, "limit", 10, 1, 100),
	}
	ponds, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.Paginated(c, "berhasil mengambil daftar wadah", ponds, response.NewMeta(f.Page, f.Limit, total))
}

func (h *PondHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "berhasil mengambil detail wadah", p)
}

func (h *PondHandler) Create(c *gin.Context) {
	var in model.PondInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	p, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.Created(c, "wadah berhasil dibuat", p)
}

func (h *PondHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.PondInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "wadah berhasil diperbarui", p)
}

func (h *PondHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err, "wadah")
		return
	}
	response.OK(c, "wadah berhasil dihapus", nil)
}

var _ = http.StatusOK // jaga import net/http bila belum terpakai
