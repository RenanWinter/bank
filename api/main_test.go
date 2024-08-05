package api

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	mockdb "github.com/RenanWinter/bank/db/mock"
	db "github.com/RenanWinter/bank/db/sqlc"
	"github.com/RenanWinter/bank/util/config"
	"github.com/RenanWinter/bank/util/random"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	cfg := config.Config{
		TokenSymmetricKey: random.String(32),
		TokenDuration:     time.Minute * 15,
	}

	server, err := NewServer(store, cfg)
	if err != nil {
		t.Fatal("cannot create server:", err)
	}
	return server
}

func newTestStructure(t *testing.T) (*gomock.Controller, *mockdb.MockStore, *Server, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()
	return ctrl, store, server, recorder
}

func TestMain(m *testing.M) {
	config.LoadConfig("../")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
