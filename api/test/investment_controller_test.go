package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	investment_model "stt/api/models/investment"
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

func TestCreateInvestmentAPI(t *testing.T) {
	investment := randomInvestment()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		params        investment_model.CreateInvestmentModel
		buildStubs    func(investmentService *mock_service.MockIInvestmentService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success OK",
			params: investment_model.CreateInvestmentModel{
				AccountID:   investment.AccountID,
				Ticker:      investment.Ticker,
				CompanyName: investment.CompanyName.String,
				MarketPrice: investment.MarketPrice,
				Description: investment.Description.String,
			},
			buildStubs: func(investmentService *mock_service.MockIInvestmentService) {
				investmentService.EXPECT().Create(gomock.Any(), gomock.Eq(db.CreateInvestmentParams{
					AccountID:   investment.AccountID,
					Ticker:      investment.Ticker,
					CompanyName: investment.CompanyName,
					Description: investment.Description,
					MarketPrice: investment.MarketPrice,
				})).Times(1).Return(investment, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	investmentService := mock_service.NewMockIInvestmentService(ctrl)

	for i := range testCases {
		tc := testCases[i]
		// run each case as a separate sub-test of this unit test
		t.Run(tc.name, func(t *testing.T) {
			// build stubs
			tc.buildStubs(investmentService)

			// start test server and send request
			server := bootstrap.NewServerApp("../..")
			route.InitInvestmentRouter(server.Engine.Group(""), investmentService)

			requestBody, err := json.Marshal(tc.params)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/investment", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.Engine.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}
}

func randomInvestment() db.Investment {
	return db.Investment{
		ID:            util.RandomInt(1, 1000),
		AccountID:     util.RandomInt(1, 1000),
		Ticker:        util.RandomString(4),
		CompanyName:   pgtype.Text{String: util.RandomString(10), Valid: true},
		BuyVolume:     int32(util.RandomInt(1, 1000)),
		BuyValue:      util.RandomInt(100, 10000),
		CapitalCost:   util.RandomInt(0, 1000),
		MarketPrice:   util.RandomMoney(),
		SellVolume:    int32(util.RandomInt(1, 1000)),
		SellValue:     util.RandomInt(1, 1000),
		CurrentVolume: int32(util.RandomInt(0, 1000)),
		Description:   pgtype.Text{String: util.RandomString(5), Valid: true},
		Status:        db.InvestmentStatusActive,
		Fee:           0,
		Tax:           0,
		UpdatedDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
}
