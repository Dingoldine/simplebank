package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  float64(util.RandomBalance()),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)

	racc, err := testQueries.GetAccount(context.Background(), acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, racc)

	require.Equal(t, racc.Owner, acc.Owner)
	require.Equal(t, racc.Balance, acc.Balance)
	require.Equal(t, racc.ID, acc.ID)
	require.Equal(t, racc.CreatedAt, acc.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)

	d := UpdateAccountParams{
		ID:      acc.ID,
		Balance: float64(util.RandomBalance()),
	}

	racc, err := testQueries.UpdateAccount(context.Background(), d)

	require.NoError(t, err)
	require.Equal(t, racc.Owner, acc.Owner)
	require.Equal(t, racc.Balance, d.Balance)
	require.Equal(t, racc.ID, acc.ID)
	require.Equal(t, racc.CreatedAt, acc.CreatedAt)

}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)

	require.NoError(t, err)

	racc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, racc)
}

func TestListAccount(t *testing.T) {

	numAcc := util.RandomInt(1, 10)

	for i := 0; i < int(numAcc); i++ {
		createRandomAccount(t)
	}

	limit := util.RandomInt(1, numAcc)

	l := ListAccountsParams{
		Limit:  int32(limit),
		Offset: int32(util.RandomInt(0, numAcc-limit)),
	}

	accs, err := testQueries.ListAccounts(context.Background(), l)
	require.NoError(t, err)
	require.NotEmpty(t, accs)
	require.Equal(t, len(accs), int(limit))
}
