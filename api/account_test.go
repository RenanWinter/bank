package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mockdb "github.com/RenanWinter/bank/db/mock"
	db "github.com/RenanWinter/bank/db/sqlc"
	"github.com/RenanWinter/bank/util/random"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()
	randomUUID, err := uuid.Parse(random.UUID())
	require.NoError(t, err)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl, store, server, recorder := newTestStructure(t)
			defer ctrl.Finish()
			tc.buildStubs(store)
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request := httptest.NewRequest("GET", url, nil)
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, randomUUID.String(), time.Minute)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:            1,
		Uuid:          uuid.New(),
		Name:          random.String(10),
		OwnerID:       1,
		AccountTypeID: 1,
		Balance:       100,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     sql.NullTime{Valid: true},
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account.ID, gotAccount.ID)
	require.Equal(t, account.Uuid, gotAccount.Uuid)
	require.Equal(t, account.Name, gotAccount.Name)
	require.Equal(t, account.OwnerID, gotAccount.OwnerID)
	require.Equal(t, account.AccountTypeID, gotAccount.AccountTypeID)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.WithinDuration(t, account.CreatedAt, gotAccount.CreatedAt, time.Second)
	require.WithinDuration(t, account.UpdatedAt, gotAccount.UpdatedAt, time.Second)
	require.Equal(t, account.DeletedAt, gotAccount.DeletedAt)
}

type eqCreateUserParamsMatcher struct {
	arg db.CreateAccountParams
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAccountParams)
	if !ok {
		return false
	}

	return arg.Name == e.arg.Name &&
		arg.OwnerID == e.arg.OwnerID &&
		arg.AccountTypeID == e.arg.AccountTypeID &&
		arg.Balance == e.arg.Balance
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %+v", e.arg)
}

func eqCreateUserParams(arg db.CreateAccountParams) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg}
}

func TestCreateAccount(t *testing.T) {
	account := randomAccount()
	randomUUID, err := uuid.Parse(random.UUID())
	require.NoError(t, err)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":            account.Name,
				"owner_id":        account.OwnerID,
				"account_type_id": account.AccountTypeID,
				"balance":         account.Balance,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "FaileCreatingAccount",
			body: gin.H{
				"name":            account.Name,
				"owner_id":        account.OwnerID,
				"account_type_id": account.AccountTypeID,
				"balance":         account.Balance,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)

				arg := db.CreateAccountParams{
					Name:          account.Name,
					OwnerID:       account.OwnerID,
					AccountTypeID: account.AccountTypeID,
					Balance:       account.Balance,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), eqCreateUserParams(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Missing Parameter",
			body: gin.H{
				"name":            account.Name,
				"owner_id":        account.OwnerID,
				"account_type_id": account.AccountTypeID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Your request is invalid")
				fmt.Println(recorder.Body)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H{
				"name":            account.Name,
				"owner_id":        account.OwnerID,
				"account_type_id": account.AccountTypeID,
				"balance":         account.Balance,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetUserByUUID(gomock.Any(), gomock.Eq(randomUUID)).
					Times(1).
					Return(db.User{ID: account.OwnerID + 1}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl, store, server, recorder := newTestStructure(t)
			defer ctrl.Finish()
			tc.buildStubs(store)
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, randomUUID.String(), time.Minute)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
