package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// userError menangani error khusus user, sisanya diteruskan ke handleError umum.
func userError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrDuplicate):
		response.Error(c, http.StatusConflict, "email sudah terdaftar")
	case errors.Is(err, service.ErrSelfAction):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		handleError(c, err, "user")
	}
}

func (h *UserHandler) List(c *gin.Context) {
	f := model.UserFilter{
		Search:   c.Query("search"),
		Role:     c.Query("role"),
		IsActive: queryBool(c, "is_active"),
		Page:     queryInt(c, "page", 1, 1, 100000),
		Limit:    queryInt(c, "limit", 10, 1, 100),
	}
	users, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		userError(c, err)
		return
	}
	response.Paginated(c, "berhasil mengambil daftar user", users, response.NewMeta(f.Page, f.Limit, total))
}

func (h *UserHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	u, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		userError(c, err)
		return
	}
	response.OK(c, "berhasil mengambil detail user", u)
}

func (h *UserHandler) Create(c *gin.Context) {
	var in model.UserCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	u, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		userError(c, err)
		return
	}
	response.Created(c, "user berhasil dibuat", u)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.UserUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	u, err := h.svc.Update(c.Request.Context(), c.GetInt64(middleware.CtxUserID), id, in)
	if err != nil {
		userError(c, err)
		return
	}
	response.OK(c, "user berhasil diperbarui", u)
}

func (h *UserHandler) SetStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.UserStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	u, err := h.svc.SetActive(c.Request.Context(), c.GetInt64(middleware.CtxUserID), id, *in.IsActive)
	if err != nil {
		userError(c, err)
		return
	}
	msg := "user berhasil dinonaktifkan"
	if *in.IsActive {
		msg = "user berhasil diaktifkan"
	}
	response.OK(c, msg, u)
}
