package school

import (
	"github.com/ujianmu/go-skeleton/internal/auth"
	"github.com/ujianmu/go-skeleton/internal/entity"
	"github.com/ujianmu/go-skeleton/internal/test"
	"github.com/ujianmu/go-skeleton/pkg/log"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.School{
		{"123", "school123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/school", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/school/123", "", nil, http.StatusOK, `*school123*`},
		{"get unknown", "GET", "/school/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/school", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/schools", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/school", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/school", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/school/123", `{"name":"schoolxyz"}`, header, http.StatusOK, "*schoolxyz*"},
		{"update verify", "GET", "/school/123", "", nil, http.StatusOK, `*schoolxyz*`},
		{"update auth error", "PUT", "/school/123", `{"name":"schoolxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/school/123", `"name":"schoolxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/school/123", ``, header, http.StatusOK, "*schoolxyz*"},
		{"delete verify", "DELETE", "/school/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/school/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
