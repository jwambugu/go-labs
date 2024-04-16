package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go-labs/internal/repository"
	"go-labs/internal/testutils"
	"go-labs/internal/testutils/factory"
	"go-labs/svc/auth"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type testSrv struct {
	db        *sqlx.DB
	mux       *http.ServeMux
	repoStore *repository.Store
}

func testServer(t *testing.T) *testSrv {
	var (
		db     = testutils.DB()
		logger = zaptest.NewLogger(t)
	)

	jwtManager, err := auth.NewPasetoToken("86f5778df1b11e35caf8bc793391bfd1")
	require.NoError(t, err)

	repoStore := repository.NewStore()
	repoStore.User = repository.NewUserRepo(db)

	var (
		authSVC = auth.NewAuthSvc(logger, repoStore, jwtManager)
		api     = NewApi(repoStore, authSVC)
	)

	return &testSrv{
		db:        db,
		mux:       api.Router(),
		repoStore: repoStore,
	}
}

func TestApi_Register_Registers(t *testing.T) {
	t.Parallel()

	var (
		user        = factory.NewUser()
		registerReq = &auth.RegisterRequest{
			Name:     user.Name,
			Email:    user.Email,
			Password: string(user.Password),
		}
	)

	b, err := json.Marshal(registerReq)
	require.NoError(t, err)
	require.NotNil(t, b)

	var (
		req = httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBuffer(b))
		rr  = httptest.NewRecorder()
	)

	srv := testServer(t)
	srv.mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp *successResponse

	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, user.Name, resp.User.Name)
	require.Equal(t, user.Email, resp.User.Email)
	require.Nil(t, resp.User.Password)
}

func TestApi_Register_FlagsDuplicateEmails(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		srv  = testServer(t)
		user = factory.NewUser()
	)

	err := srv.repoStore.User.Create(ctx, user)
	require.NoError(t, err)

	registerReq := &auth.RegisterRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: "password",
	}

	b, err := json.Marshal(registerReq)
	require.NoError(t, err)
	require.NotNil(t, b)

	var (
		req = httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBuffer(b))
		rr  = httptest.NewRecorder()
	)

	srv.mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestApi_Register_Validates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		req         func() *auth.RegisterRequest
		errorsCount int
	}{
		{
			req:         func() *auth.RegisterRequest { return nil },
			errorsCount: 3,
		},
		{
			req: func() *auth.RegisterRequest {
				user := factory.NewUser()
				return &auth.RegisterRequest{Name: user.Name}
			},
			errorsCount: 2,
		},
		{
			req: func() *auth.RegisterRequest {
				user := factory.NewUser()
				return &auth.RegisterRequest{
					Name:  user.Name,
					Email: user.Email,
				}
			},
			errorsCount: 1,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tt.req())
			require.NoError(t, err)
			require.NotNil(t, b)

			var (
				req = httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBuffer(b))
				rr  = httptest.NewRecorder()
			)

			srv := testServer(t)
			srv.mux.ServeHTTP(rr, req)

			require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

			var resp *errorResponse

			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)
			require.Len(t, resp.Errors, tt.errorsCount)
		})
	}
}
