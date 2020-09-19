package genity

import (
	"net/http"
	"testing"
	"time"

	"github.com/wondenge/go-genity/internal/auth"
	"github.com/wondenge/go-genity/internal/entity"
	"github.com/wondenge/go-genity/internal/test"
	"github.com/wondenge/go-genity/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Genity{
		{"123", "genity123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/genitys", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/genitys/123", "", nil, http.StatusOK, `*genity123*`},
		{"get unknown", "GET", "/genitys/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/genitys", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/genitys", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/genitys", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/genitys", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/genitys/123", `{"name":"genityxyz"}`, header, http.StatusOK, "*genityxyz*"},
		{"update verify", "GET", "/genitys/123", "", nil, http.StatusOK, `*genityxyz*`},
		{"update auth error", "PUT", "/genitys/123", `{"name":"genityxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/genitys/123", `"name":"genityxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/genitys/123", ``, header, http.StatusOK, "*genityxyz*"},
		{"delete verify", "DELETE", "/genitys/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/genitys/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
