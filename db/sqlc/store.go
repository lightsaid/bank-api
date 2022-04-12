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

	// err := store.execTx(ctx, func(q *Queries)error {
	// 	var err error

	// 	result.Transfer, err := q.Crea
	// })

	return result, nil
}
