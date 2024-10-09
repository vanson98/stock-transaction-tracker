package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// service mock
	accountService := mock_service.NewMockIAccountService(ctrl)

	// build stubs
	accountService.EXPECT().CreateNew(gomock.Any(), gomock.Eq(db.CreateAccountParams{
		ChannelName: account.ChannelName,
		Owner:       account.Owner,
		Currency:    account.Currency,
	})).Times(1).Return(account, nil)

	// start test server and send request
	server := bootstrap.NewServerApp("../..")
	routeg := server.Engine.Group("")
	route.InitAccountRouter(routeg, accountService)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account")
	requestBody, err := json.Marshal(apimodels.CreateAccountRequest{
		ChannelName: account.ChannelName,
		Owner:       account.Owner,
		Currency:    account.Currency,
	})
	require.NoError(t, err)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	require.NoError(t, err)

	server.Engine.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomAccount() db.Account {
	return db.Account{
		ID:          util.RandomInt(1, 1000),
		Owner:       util.RandomOwner(),
		Balance:     0,
		ChannelName: util.RandomString(3),
		Currency:    util.RandomCurrency(),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
}
