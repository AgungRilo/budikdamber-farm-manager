package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		role    string
		allowed []string
		want    int
	}{
		{"owner boleh akses route owner", RoleOwner, []string{RoleOwner}, http.StatusOK},
		{"worker ditolak di route owner", RoleWorker, []string{RoleOwner}, http.StatusForbidden},
		{"tanpa role ditolak", "", []string{RoleOwner}, http.StatusForbidden},
		{"worker boleh di route bersama", RoleWorker, []string{RoleOwner, RoleWorker}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/x",
				func(c *gin.Context) { // simulasi AuthRequired
					if tt.role != "" {
						c.Set(CtxRole, tt.role)
					}
				},
				RequireRoles(tt.allowed...),
				func(c *gin.Context) { c.Status(http.StatusOK) },
			)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))

			if w.Code != tt.want {
				t.Errorf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}
