package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 提供一个执行所有事务需要 db 和 Queries
type Store struct {
	db      *sql.DB
	Queries *Queries
}

// NewStore 创建一个 Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx 执行数据库的 transaction
// fn Queries struct 的方法，也就是说 fn 就是 crud 的方法
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 事务开始
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		// 如果有错，就是执行事务回滚
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// 提交事务
	return tx.Commit()
}

// 下面 定义业务的事务处理函数

// TransferTxParams 转账记录参数定义
type TransferTxParams struct {
	FromAccountID int64 `josn:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferResult 事务成功结果集
type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx 转账事务处理
// 1. 在transfer表创建转账记录 2. 往entries添加条目 3. 更新账户 balance
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error) {
	var result TransferResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// 转账记录
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// 创建转账人 entries 条目数(流水账)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // 负数，转账
		})

		if err != nil {
			return err
		}

		// 创建 entries 条目数(流水账)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // 整数，入账
		})
		if err != nil {
			return err
		}

		//  NOTE: 通过如果这里没有锁机制，一旦发生并发，就会错误 （事务隔离）
		account1, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount, err = q.UpdateAcount(ctx, UpdateAcountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		result.ToAccount, err = q.UpdateAcount(ctx, UpdateAcountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
