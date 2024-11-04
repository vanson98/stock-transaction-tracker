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
	"stt/services/dtos"
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
				ChannelName: "TCB",
				Owner:       account.Owner,
				Currency:    "account.Currency",
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

			request, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.Engine.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}
}

func TestGetAccountByIdAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountService := mock_service.NewMockIAccountService(ctrl)
	account := randomAccount()

	testCases := []struct {
		name          string
		param         int64
		buildStubs    func(accountService *mock_service.MockIAccountService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: account.ID,
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().GetById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Bad Request",
			param: -1,
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().GetById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal Server Error",
			param: account.ID,
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().GetById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/accounts/%d", subTest.param)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)
			responseRecorder := httptest.NewRecorder()

			server.Engine.ServeHTTP(responseRecorder, request)

			subTest.checkResponse(t, responseRecorder)
		})

	}
}

func TestTranserMoneyAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	account := randomAccount()
	account.CreatedAt = pgtype.Timestamptz{}

	entry := randomEntry()
	entry.CreatedAt = pgtype.Timestamptz{}

	entry.AccountID = account.ID
	amount := util.RandomInt(-100, 100)
	transferResultDto := randomTransferResultDto(account, entry, amount)

	accountService := mock_service.NewMockIAccountService(ctrl)

	testCases := []struct {
		name          string
		param         apimodels.TransferMoneyRequest
		buildStubs    func(accService *mock_service.MockIAccountService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success Transfer",
			param: apimodels.TransferMoneyRequest{
				AccountID: account.ID,
				Amount:    amount,
				EntryType: string(entry.Type),
				Currency:  account.Currency,
			},
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().TransferMoney(gomock.Any(), dtos.TransferMoneyTxParam{
					AccountID: account.ID,
					Amount:    amount,
					EntryType: "IT",
				}).Times(1).Return(transferResultDto, nil)
				accService.EXPECT().GetById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccountTranserResult(t, recorder.Body, transferResultDto)
			},
		},
		{
			name: "Bad request",
			param: apimodels.TransferMoneyRequest{
				AccountID: -2,
				Amount:    amount,
				EntryType: string(entry.Type),
				Currency:  "TEST",
			},
			buildStubs: func(accService *mock_service.MockIAccountService) {
				accService.EXPECT().TransferMoney(gomock.Any(), gomock.Any()).Times(0)
				accService.EXPECT().GetById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		subTest := testCases[i]
		t.Run(subTest.name, func(t *testing.T) {
			subTest.buildStubs(accountService)

			server := bootstrap.NewServerApp("../..")
			route.InitAccountRouter(server.Engine.Group(""), accountService)
			body, err := json.Marshal(subTest.param)
			require.NoError(t, err)

			request, err := http.NewRequest("PUT", "/account-transfer", bytes.NewBuffer(body))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.Engine.ServeHTTP(recorder, request)

			subTest.checkResponse(t, recorder)
		})
	}
}

func TestGetAllAccountAPI(t *testing.T) {
	// create 2 account for testing
	var accounts []db.Account
	accounts = append(accounts, randomAccount())
	accounts = append(accounts, randomAccount())

	ctrl := gomock.NewController(t)
	accountService := mock_service.NewMockIAccountService(ctrl)

	testCase := []struct {
		name          string
		buildStub     func(param *mock_service.MockIAccountService)
		checkResponse func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name: "Happy Case",
			buildStub: func(mockService *mock_service.MockIAccountService) {
				mockService.EXPECT().ListAllAccount(gomock.Any()).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, response.Result().StatusCode, 200)
				var responseDataModel []db.Account

				data, _ := io.ReadAll(response.Body)
				err := json.Unmarshal(data, &responseDataModel)
				require.NoError(t, err)

				require.Len(t, responseDataModel, 2)
				require.Greater(t, responseDataModel[0].ID, int64(0))
				require.Greater(t, responseDataModel[1].ID, int64(0))
			},
		},
	}
	// sub test
	for _, subTest := range testCase {
		t.Run(subTest.name, func(t *testing.T) {
			subTest.buildStub(accountService)

			server := bootstrap.NewServerApp("../..")
			route.InitAccountRouter(&server.Engine.RouterGroup, accountService)
			httpRequest, err := http.NewRequest("GET", "/accounts", nil)
			require.NoError(t, err)

			httpResponse := httptest.NewRecorder()
			server.Engine.ServeHTTP(httpResponse, httpRequest)

			subTest.checkResponse(t, httpResponse)
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

func randomEntry() db.Entry {
	return db.Entry{
		ID:        util.RandomInt(1, 100),
		AccountID: util.RandomInt(1, 200),
		Amount:    util.RandomInt(-100, 100),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		Type: db.EntryType(util.RandomEntryType()),
	}
}

func randomTransferResultDto(acc db.Account, entry db.Entry, amount int64) dtos.TransferMoneyTxResult {
	acc.Balance += amount
	entry.Amount = amount
	return dtos.TransferMoneyTxResult{
		UpdatedAccount: acc,
		Entry:          entry,
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

func requireBodyMatchAccountTranserResult(t *testing.T, body *bytes.Buffer, transferResult dtos.TransferMoneyTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	transferRes := dtos.TransferMoneyTxResult{}

	err = json.Unmarshal(data, &transferRes)
	require.NoError(t, err)

	require.Equal(t, transferRes.Entry.AccountID, transferRes.UpdatedAccount.ID)
	require.Equal(t, transferResult.UpdatedAccount, transferRes.UpdatedAccount)
	require.Equal(t, transferResult.Entry, transferRes.Entry)
}
