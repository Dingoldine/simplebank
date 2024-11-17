package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)

	if err != nil {
		return err
	}
	q := New(tx)

	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, t TransferTxParams) (TransferTxResult, error) {
	var r TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		r.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(t))
		if err != nil {
			return err
		}

		r.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: t.FromAccountID,
			Amount:    -t.Amount,
		})
		if err != nil {
			return err
		}

		r.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: t.ToAccountID,
			Amount:    t.Amount,
		})
		if err != nil {
			return err
		}

		r.FromAccount, err = q.GetAccountForUpdate(ctx, t.FromAccountID)
		if err != nil {
			return err
		}

		r.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      t.FromAccountID,
			Balance: r.FromAccount.Balance - t.Amount,
		})
		if err != nil {
			return err
		}

		r.ToAccount, err = q.GetAccountForUpdate(ctx, t.ToAccountID)
		if err != nil {
			return err
		}

		r.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      t.ToAccountID,
			Balance: r.ToAccount.Balance + t.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return r, err

}
