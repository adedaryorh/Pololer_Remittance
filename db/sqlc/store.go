package db

import (
	"context"
	"database/sql"
	"fmt"
)

//#begin trans
//transfer money
//deposit
//withrawal
//update balance
//#commit trans

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fq func(q *Queries) error) error {
	//Initialized Trans
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fq(q)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("Encountered roolback error: %v", txErr)
		}
		return err
	}
	return tx.Commit()
}
