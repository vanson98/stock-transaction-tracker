package service_test

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createDefaultInvestmnet(t *testing.T, accountId int64) db.Investment {
	arg := db.CreateInvestmentParams{
		AccountID:     accountId,
		Ticker:        util.RandomUpperString(3),
		CompanyName:   pgtype.Text{},
		BuyVolume:     0,
		BuyValue:      0,
		CapitalCost:   0,
		MarketPrice:   0,
		SellVolume:    0,
		SellValue:     0,
		CurrentVolume: 0,
		Description:   pgtype.Text{},
		Status:        db.InvestmentStatusInactive,
		Fee:           0,
		Tax:           0,
	}
	ivm, err := investmentService.Create(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ivm)
	require.NotZero(t, ivm.ID)
	require.Equal(t, accountId, ivm.AccountID)
	require.Equal(t, arg.Ticker, ivm.Ticker)
	require.Equal(t, arg.BuyValue, ivm.BuyValue)
	require.Equal(t, arg.BuyVolume, ivm.BuyVolume)
	require.Equal(t, arg.CapitalCost, ivm.CapitalCost)
	require.Equal(t, arg.SellVolume, ivm.SellVolume)
	return ivm
}

func TestCreateInvestment(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	createDefaultInvestmnet(t, acc.ID)
}

func TestGetInvestmentById(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	ivm := createDefaultInvestmnet(t, acc.ID)

	investment, err := investmentService.GetById(context.Background(), ivm.ID)
	require.NoError(t, err)
	require.NotEmpty(t, investment)
	require.Equal(t, ivm.AccountID, investment.AccountID)
	require.Equal(t, ivm.Ticker, investment.Ticker)
	require.Equal(t, ivm.BuyValue, investment.BuyValue)
	require.Equal(t, ivm.BuyVolume, investment.BuyVolume)
	require.Equal(t, ivm.CapitalCost, investment.CapitalCost)
	require.Equal(t, ivm.SellVolume, investment.SellVolume)
	require.Equal(t, ivm.Status, investment.Status)
	require.Equal(t, ivm.Fee, investment.Fee)
	require.Equal(t, ivm.Tax, investment.Tax)
	require.Equal(t, ivm.UpdatedDate, investment.UpdatedDate)
	require.Equal(t, ivm.CurrentVolume, investment.CurrentVolume)
	require.Equal(t, ivm.Description, investment.Description)
	require.Equal(t, ivm.CompanyName, investment.CompanyName)
	require.Equal(t, ivm.MarketPrice, investment.MarketPrice)
	require.Equal(t, ivm.SellValue, investment.SellValue)
	require.Equal(t, ivm.ID, investment.ID)
}

func TestSearchInvestment(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	n := 10
	for i := 0; i < n; i++ {
		createDefaultInvestmnet(t, acc.ID)
	}

	investments, err := investmentService.SearchPaging(context.Background(), db.SearchInvestmentPagingParams{
		AccountIds: []int64{acc.ID},
		SearchText: "",
		OrderBy:    "ticker",
		SortType:   "descending",
		FromOffset: 0,
		TakeLimit:  10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, investments)
	require.Len(t, investments, n)
}

func TestUpdateMarketPrice(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	ivm := createDefaultInvestmnet(t, acc.ID)
	updateMarketPrice := util.RandomInt(1000, 100000)
	result, err := investmentService.UpdateMarketPrice(context.Background(), db.UpdateMarketPriceParams{
		ID:          ivm.ID,
		MarketPrice: updateMarketPrice,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, updateMarketPrice, result.MarketPrice)
	require.Equal(t, ivm.ID, result.ID)
}
