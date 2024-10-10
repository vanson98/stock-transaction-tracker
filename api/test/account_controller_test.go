package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	apimodels "stt/api/models"
	"stt/api/route"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	mock_service "stt/services/mock"
	"stt/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	account.Balance = 0
	account.CreatedAt = pgtype.Timestamptz{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		params        apimodels.CreateAccountRequest
		buildStubs    func(accountService *mock_service.MockIAccountService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success OK",
			params: apimodels.CreateAccountRequest{
				ChannelName: account.ChannelName,
				Owner:       account.Owner,
				Currency:    account.Currency,
			},
			buildStubs: func(accountService *mock_service.MockIAccountService) {
				accountService.EXPECT().CreateNew(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					ChannelName: account.ChannelName,
					Owner:       account.Owner,
					Currency:    account.Currency,
				})).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "Bad Request",
			params: apimodels.CreateAccountRequest{
				ChannelName: "",
				Owner:       account.Owner,
				Currency:    account.Currency,
			},
			buildStubs: func(accountService *mock_service.MockIAccountService) {
				accountService.EXPECT().CreateNew(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			params: apimodels.CreateAccountRequest{
				ChannelName: account.ChannelName,
				Owner:       account.Owner,
				Currency:    account.Currency,
			},
			buildStubs: func(accountService *mock_service.MockIAccountService) {
				accountService.EXPECT().CreateNew(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					ChannelName: account.ChannelName,
					Owner:       account.Owner,
					Balance:     account.Balance,
					Currency:    account.Currency,
				})).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	accountService := mock_service.NewMockIAccountService(ctrl)

	for i := range testCases {
		tc := testCases[i]
		// run each case as a separate sub-test of this unit test
		t.Run(tc.name, func(t *testing.T) {
			// build stubs
			tc.buildStubs(accountService)

			// start test server and send request
			server := bootstrap.NewServerApp("../..")
			route.InitAccountRouter(server.Engine.Group(""), accountService)

			requestBody, err := json.Marshal(tc.params)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/account", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.Engine.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}
}

func TestGetAccountAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountService := mock_service.NewMockIAccountService(ctrl)
	account := randomAccount()

	testCases := []struct {
		name          string
		param         apimodels.GetAccountRequest
		buildStubs    func(accountService *mock_service.MockIAccountService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: apimodels.GetAccountRequest{Id: account.ID},
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().GetById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		subTest := testCases[i]

		t.Run(subTest.name, func(t *testing.T) {
			subTest.buildStubs(accountService)

			// create demo server
			server := bootstrap.NewServerApp("../..")
			route.InitAccountRouter(server.Engine.Group(""), accountService)

			//body, err := json.Marshal(subTest.param)
			//require.NoError(t, err)

			url := fmt.Sprintf("/account/%d", account.ID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)
			responseRecorder := httptest.NewRecorder()

			server.Engine.ServeHTTP(responseRecorder, request)

			subTest.checkResponse(t, responseRecorder)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:          util.RandomInt(1, 1000),
		Owner:       util.RandomOwner(),
		Balance:     util.RandomMoney(),
		ChannelName: util.RandomString(3),
		Currency:    util.RandomCurrency(),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
