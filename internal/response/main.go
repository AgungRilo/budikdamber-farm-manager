package response

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Body{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, Body{Success: true, Message: message, Data: data})
}

func Paginated(c *gin.Context, message string, data any, meta Meta) {
	c.JSON(http.StatusOK, Body{Success: true, Message: message, Data: data, Meta: &meta})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Body{Success: false, Message: message})
}

// Abort dipakai di middleware: kirim error dan hentikan handler berikutnya.
func Abort(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Body{Success: false, Message: message})
}

func NewMeta(page, limit int, total int64) Meta {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return Meta{Page: page, Limit: limit, Total: total, TotalPages: totalPages}
}

// ValidationError mengubah error binding Gin menjadi pesan per field.
func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			fields[fe.Field()] = fieldMessage(fe)
		}
		c.JSON(http.StatusBadRequest, Body{Success: false, Message: "validasi gagal", Errors: fields})
		return
	}
	c.JSON(http.StatusBadRequest, Body{Success: false, Message: "format request tidak valid"})
}

// UseJSONFieldNames membuat nama field di pesan error memakai tag json
// ("email"), bukan nama struct Go ("Email"). Panggil sekali saat startup.
func UseJSONFieldNames() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("minimal %s karakter", fe.Param())
		}
		return fmt.Sprintf("minimal %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("maksimal %s karakter", fe.Param())
		}
		return fmt.Sprintf("maksimal %s", fe.Param())
	case "gt":
		return fmt.Sprintf("harus lebih dari %s", fe.Param())
	case "gte":
		return fmt.Sprintf("minimal %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", strings.ReplaceAll(fe.Param(), " ", ", "))
	default:
		return "tidak valid"
	}
}
