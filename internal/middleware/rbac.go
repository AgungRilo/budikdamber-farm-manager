package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
)

const (
	RoleOwner  = "owner"
	RoleWorker = "worker"
)

// RequireRoles hanya meloloskan user dengan salah satu role yang disebut.
// Wajib dipasang SETELAH AuthRequired (yang mengisi role ke context).
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		if !allowed[c.GetString(CtxRole)] {
			response.Abort(c, http.StatusForbidden, "anda tidak punya akses ke fitur ini")
			return
		}
		c.Next()
	}
}
