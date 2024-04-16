package db_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"testing"
	"time"
)

//var testQuery *db.Store

func createRandomAccount(customer_id int64, t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		CustomerID:    int32(customer_id),
		Balance:       500,
		Currency:      "NGN",
		AccountType:   "Individual",
		AccountStatus: "Active",
	}

	account, err := testQuery.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, account.Balance, arg.Balance)
	assert.Equal(t, account.Currency, arg.Currency)
	assert.Equal(t, account.AccountType, arg.AccountType)
	assert.Equal(t, account.AccountStatus, arg.AccountStatus)
	assert.Equal(t, account.CustomerID, arg.CustomerID)
	assert.WithinDuration(t, account.CreatedAt, time.Now(), 2*time.Second)

	return account
}

func TestTransfer(t *testing.T) {
	customer1 := createRandomCustomer(t)
	customer2 := createRandomCustomer(t)

	account1 := createRandomAccount(customer1.ID, t)
	account2 := createRandomAccount(customer2.ID, t)

	arg := db.CreateTransferParams{
		FromAccountID: int32(account1.ID),
		ToAccountID:   int32(account2.ID),
		Amount:        10,
	}

	txResponseChan := make(chan db.TransferTxResponse)
	errorchan := make(chan error)

	count := 3

	for i := 0; i < 3; i++ {
		go func() {
			tx, err := testQuery.TransferTx(context.Background(), arg)
			errorchan <- err
			txResponseChan <- tx
		}()
	}

	for i := 0; i < count; i++ {
		err := <-errorchan
		tx := <-txResponseChan

		assert.NoError(t, err)
		assert.NotEmpty(t, tx)
		//test transfer
		assert.Equal(t, tx.Transfer.FromAccountID, arg.FromAccountID)
		assert.Equal(t, tx.Transfer.ToAccountID, arg.ToAccountID)
		assert.Equal(t, tx.Transfer.Amount, arg.Amount)
		//test entry
		//entry IN
		assert.Equal(t, tx.EntryIn.AccountID, arg.ToAccountID)
		assert.Equal(t, tx.EntryIn.Amount, arg.Amount)
		//entry out
		assert.Equal(t, tx.EntryOut.AccountID, arg.FromAccountID)
		assert.Equal(t, tx.EntryOut.Amount, -1*arg.Amount)
	}

	newAccount1, err := testQuery.GetAccountByID(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount1)
	newAccount2, err := testQuery.GetAccountByID(context.Background(), account2.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount2)

	newAccount := float64(count * int(arg.Amount))
	assert.Equal(t, newAccount1.Balance, account1.Balance-newAccount)
	assert.Equal(t, newAccount2.Balance, account1.Balance+newAccount)

}
