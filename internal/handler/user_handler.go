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

// List godoc
// @Summary   Daftar user (owner)
// @Tags      Users
// @Produce   json
// @Security  BearerAuth
// @Param     page    query  int     false  "Halaman"  default(1)
// @Param     limit   query  int     false  "Jumlah per halaman"  default(10)
// @Param     search  query  string  false  "Cari nama / email"
// @Success   200  {object}  response.Body{data=[]model.User,meta=response.Meta}
// @Failure   403  {object}  response.Body
// @Router    /users [get]
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

// Get godoc
// @Summary   Detail user (owner)
// @Tags      Users
// @Produce   json
// @Security  BearerAuth
// @Param     id   path      int  true  "ID user"
// @Success   200  {object}  response.Body{data=model.User}
// @Failure   404  {object}  response.Body
// @Router    /users/{id} [get]
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

// Create godoc
// @Summary   Buat user (owner)
// @Tags      Users
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     body  body      model.UserCreateInput  true  "Data user"
// @Success   201   {object}  response.Body{data=model.User}
// @Failure   400   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /users [post]
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

// Update godoc
// @Summary   Ubah user (owner)
// @Tags      Users
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                    true  "ID user"
// @Param     body  body      model.UserUpdateInput  true  "Data user"
// @Success   200   {object}  response.Body{data=model.User}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Router    /users/{id} [put]
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

// SetStatus godoc
// @Summary   Aktifkan / nonaktifkan user (owner)
// @Tags      Users
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     id    path      int                    true  "ID user"
// @Param     body  body      model.UserStatusInput  true  "Status"
// @Success   200   {object}  response.Body{data=model.User}
// @Failure   400   {object}  response.Body
// @Failure   404   {object}  response.Body
// @Router    /users/{id}/status [patch]
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

// Register godoc
// @Summary   Daftar akun baru (publik, menunggu persetujuan owner)
// @Tags      Auth
// @Accept    json
// @Produce   json
// @Param     body  body      model.RegisterInput  true  "Data pendaftaran"
// @Success   201   {object}  response.Body{data=model.User}
// @Failure   400   {object}  response.Body
// @Failure   409   {object}  response.Body
// @Router    /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var in model.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.ValidationError(c, err)
		return
	}
	u, err := h.svc.Register(c.Request.Context(), in)
	if err != nil {
		userError(c, err)
		return
	}
	response.Created(c, "pendaftaran berhasil, akun menunggu aktivasi dari owner", u)
}
