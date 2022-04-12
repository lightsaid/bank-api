package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"lightsaid.com/build-api/bank-api/util"
)

func createRandomCntry(t *testing.T, account *Account) Entry {
	params := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandromAccount(t)
	_ = createRandomCntry(t, &account)

}

func TestGetEntry(t *testing.T) {
	account := createRandromAccount(t)
	entry1 := createRandomCntry(t, &account)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandromAccount(t)

	for i := 0; i < 10; i++ {
		createRandomCntry(t, &account)
	}

	arg := ListEntriesParams{
		Limit:     5,
		Offset:    5,
		AccountID: account.ID,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
