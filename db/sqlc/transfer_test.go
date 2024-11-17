package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransaction(t *testing.T) Transfer {
	s := createRandomAccount(t)
	r := createRandomAccount(t)

	p := CreateTransferParams{
		FromAccountID: s.ID,
		ToAccountID:   r.ID,
		Amount:        float64(util.RandomBalance()),
	}

	trans, err := testQueries.CreateTransfer(context.Background(), p)

	require.NoError(t, err)
	require.NotEmpty(t, trans)
	require.Equal(t, p.FromAccountID, trans.FromAccountID)
	require.Equal(t, p.ToAccountID, trans.ToAccountID)

	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreatedAt)

	return trans
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransaction(t)
}

func TestGetTransfer(t *testing.T) {
	trans := createRandomTransaction(t)
	trans2, err := testQueries.GetTransfer(context.Background(), trans.ID)

	require.NoError(t, err)
	require.NotEmpty(t, trans)

	require.Equal(t, trans2.ID, trans.ID)
	require.Equal(t, trans2.FromAccountID, trans.FromAccountID)
	require.Equal(t, trans2.ToAccountID, trans.ToAccountID)
	require.Equal(t, trans2.CreatedAt, trans.CreatedAt)
}
