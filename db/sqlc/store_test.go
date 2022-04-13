package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// 创建两个用于转账的账号
	account1 := createRandromAccount(t)
	account2 := createRandromAccount(t)

	/*
	* NOTE: 使用事务操作需要非常小心，很容易编写，但是也容器产生噩梦,
	* 确保事务运行的最佳方式就是使用 Goroutine
	 */

	store := NewStore(testDB)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferResult)

	// 测试多条转账记录，并验证
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// 验证
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		// 检查 transfer 表记录
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// 获取 transfer 记录检查
		_, err = store.Queries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// 检查 entries 表记录 （出帐 账号）
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		// 检查 entries 表记录 （入帐 账号）
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		// TODO: 待验证 账号余额 （account 表中 balance）

	}

}
