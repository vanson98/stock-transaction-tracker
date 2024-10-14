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

	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newuser := createRandomUser(t)
	user, err := userService.GetByUserName(context.Background(), newuser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, newuser.Username, user.Username)
	require.Equal(t, newuser.Email, user.Email)
	require.Equal(t, newuser.FullName, user.FullName)
	require.Equal(t, newuser.HashedPassword, user.HashedPassword)
	require.Equal(t, newuser.PasswordChangedAt, user.PasswordChangedAt)
	require.Equal(t, newuser.CreatedAt, user.CreatedAt)
}
