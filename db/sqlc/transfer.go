package db

import "context"

type TransferTxResponse struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	EntryIn     Entry    `json:"entry_in"`
	EntryOut    Entry    `json:"entry_out"`
	Transfer    Transfer `json:"transfer"`
}

func (s *Store) TransferTx(ctx context.Context, tr CreateTransferParams) (TransferTxResponse, error) {
	var tx TransferTxResponse
	var errT error

	err := s.execTx(ctx, func(q *Queries) error {
		//transfer moni
		tx.Transfer, errT = q.CreateTransfer(context.Background(), tr)
		if errT != nil {
			return errT
		}
		//record entry
		inEntryArg := CreateEntryParams{
			AccountID: tr.ToAccountID,
			Amount:    tr.Amount,
		}
		tx.EntryIn, errT = q.CreateEntry(context.Background(), inEntryArg)
		if errT != nil {
			return errT
		}

		outEntryArg := CreateEntryParams{
			AccountID: tr.ToAccountID,
			Amount:    -1 * tr.Amount,
		}
		tx.EntryOut, errT = q.CreateEntry(context.Background(), outEntryArg)
		if errT != nil {
			return errT
		}
		//update bal of new acct
		toArg := UpdateAccountBalanceManualParams{
			Amount: tr.Amount,
			ID:     int64(tr.ToAccountID),
		}
		tx.ToAccount, errT = q.UpdateAccountBalanceManual(context.Background(), toArg)
		if errT != nil {
			return errT
		}

		fromArg := UpdateAccountBalanceManualParams{
			Amount: -1 * tr.Amount,
			ID:     int64(tr.FromAccountID),
		}
		tx.FromAccount, errT = q.UpdateAccountBalanceManual(context.Background(), fromArg)
		if errT != nil {
			return errT
		}
		return nil
	})
	return tx, err
}
