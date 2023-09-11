package db

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// FinalFunc is the final function for transaction
type FinalFunc func(error) error

// Context is the context for store
type Context struct {
	context.Context
	DB boil.ContextExecutor
}

// FromContext create new store context
func FromContext(ctx context.Context) Context {
	return Context{
		Context: ctx,
		DB:      GetDB(),
	}
}

// NewTransaction create new transaction
func NewTransaction(ctx context.Context) (Context, FinalFunc) {
	tx, err := GetDB().BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	return Context{
			Context: ctx,
			DB:      tx,
		}, func(err error) error {
			if err != nil {
				txErr := tx.Rollback()
				if txErr != nil {
					return txErr
				}
				return err
			}
			return tx.Commit()
		}
}

// Transaction create new transaction with callback function txFunc
func Transaction(ctx context.Context, txFunc func(ctx Context) error) error {
	tx, err := GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = txFunc(Context{
		Context: ctx,
		DB:      tx,
	})

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return txErr
		}
		return err
	}

	return tx.Commit()
}
