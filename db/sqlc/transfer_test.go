package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"lightsaid.com/build-api/bank-api/util"
)

func createRandomTransfer(t *testing.T, account1, account2 *Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NotEmpty(t, transfer)
	require.NoError(t, err)

	require.Equal(t, arg.FromAccountID, account1.ID)
	require.Equal(t, arg.ToAccountID, account2.ID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandromAccount(t)
	account2 := createRandromAccount(t)
	createRandomTransfer(t, &account1, &account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandromAccount(t)
	account2 := createRandromAccount(t)
	tfr := createRandomTransfer(t, &account1, &account2)

	transfer, err := testQueries.GetTransfer(context.Background(), tfr.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, tfr.Amount, transfer.Amount)

	require.WithinDuration(t, transfer.CreatedAt, tfr.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandromAccount(t)
	account2 := createRandromAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, &account1, &account2)
		createRandomTransfer(t, &account2, &account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
