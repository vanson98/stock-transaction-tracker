package service_test

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomInvestmnet(t *testing.T, accountId int64) db.Investment {
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
	acc := createRandomAccount(t)
	createRandomInvestmnet(t, acc.ID)
}
