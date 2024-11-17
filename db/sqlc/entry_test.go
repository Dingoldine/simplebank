package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(acc Account, t *testing.T) Entry {

	p := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    float64(util.RandomBalance()),
	}

	ent, err := testQueries.CreateEntry(context.Background(), p)

	require.NoError(t, err)
	require.NotEmpty(t, ent)
	require.Equal(t, p.AccountID, ent.AccountID)
	require.Equal(t, p.Amount, ent.Amount)

	require.NotZero(t, ent.ID)
	require.NotZero(t, ent.CreatedAt)

	return ent
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(acc, t)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	ent := createRandomEntry(acc, t)
	ent2, err := testQueries.GetEntry(context.Background(), ent.ID)

	require.NoError(t, err)
	require.NotEmpty(t, ent)

	require.Equal(t, ent2.ID, ent.ID)
	require.Equal(t, ent2.AccountID, ent.AccountID)
	require.Equal(t, ent2.CreatedAt, ent.CreatedAt)
	require.Equal(t, ent2.Amount, ent.Amount)
}
