package service_test

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.User {
	arg := db.CreateUserParams{
		Username:       util.RandomString(6),
		HashedPassword: "hash",
		FullName:       util.RandomString(6),
		Email:          util.RandomEmail(),
	}
	user, err := userService.CreateNew(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.NotEmpty(t, arg.FullName, user.FullName)
	require.NotEmpty(t, arg.HashedPassword, user.HashedPassword)
	require.NotEmpty(t, arg.Email, user.Email)

	//require.Equal(t, pgtype.Timestamptz{}, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
