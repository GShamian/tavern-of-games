package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GShamian/tavern-of-games/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandlerUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	s := newServer(teststore.New())
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
