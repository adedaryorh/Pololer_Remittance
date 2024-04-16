package api

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"github/adedaryorh/pooler_Remmitance_Application/utils"
	"net/http"
)

type Account struct {
	server *Server
}

func (a Account) router(server *Server) {
	a.server = server
	//AuthenticatedMiddleware()

	serverGroup := server.router.Group("/account")
	serverGroup.POST("create-account", a.createAccount)
	serverGroup.GET("", a.getCustomerAccounts)
	serverGroup.POST("transfer", a.transfer)
	serverGroup.POST("add-money", a.addMoney)
	serverGroup.POST("withdraw-money", a.withdrawMoney)

}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (a *Account) getCustomerAccounts(ctx *gin.Context) {
	customerId, err := utils.GetActiveCustomer(ctx)
	if err != nil {
		return
	}

	accounts, err := a.server.queries.GetAccountByCustomerId(context.Background(), int32(customerId))
	if err != nil {
		ctx.JSON(http.StatusOK, accounts)
	}
}

type TransferRequest struct {
	ToAccountID   int32   `json:"to-account-id" binding:"required"`
	FromAccountID int32   `json:"from_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

func (a *Account) transfer(ctx *gin.Context) {
	customerId, err := utils.GetActiveCustomer(ctx)
	tr := new(TransferRequest)

	if err := ctx.ShouldBindJSON(&tr); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.FromAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot retrieve account"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//confirm acct belong to customer
	if account.CustomerID != int32(customerId) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toAccount, err := a.server.queries.GetAccountByID(context.Background(), int64(tr.ToAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot find account to deposit"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//confirm to-acc currency  = from-acc currency
	if toAccount.Currency != account.Currency {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Currencies not match issue encoutered"})
		return
	}

	if account.Balance < tr.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Low Balance available"})
		return
	}

	txArg := db.CreateTransferParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Amount:        tr.Amount,
	}
	tx, err := a.server.queries.TransferTx(context.Background(), txArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, tx)
}

func (a Account) createAccount(ctx *gin.Context) {
	customerId, err := utils.GetActiveCustomer(ctx)
	if err != nil {
		return
	}

	acct := new(AccountRequest)
	if err := ctx.ShouldBindJSON(&acct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	argument := db.CreateAccountParams{
		CustomerID:    int32(customerId),
		Currency:      acct.Currency,
		AccountStatus: "Active",
		AccountType:   "Individual",
		Balance:       0,
	}

	account, err := a.server.queries.CreateAccount(context.Background(), argument)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "An account with currency exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	accountNumber, err := utils.GenerateAccountNumber(int64(account.ID), account.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		err := a.server.queries.DeleteAccount(context.Background(), account.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	account, err = a.server.queries.UpdateAccountNumber(context.Background(), db.UpdateAccountNumberParams{
		AccountNumber: sql.NullString{String: accountNumber, Valid: true},
		ID:            account.ID,
	})
	ctx.JSON(http.StatusCreated, account)
}

type AddMoneyRequest struct {
	ToAccountID int64   `json:"to_account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Reference   string  `json:"reference" binding:"required"`
}

func (a *Account) addMoney(ctx *gin.Context) {
	customerId, err := utils.GetActiveCustomer(ctx)
	if err != nil {
		return
	}

	obj := new(AddMoneyRequest)
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check if acct reciving exits
	account, err := a.server.queries.GetAccountByID(context.Background(), obj.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if account.CustomerID != int32(customerId) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized to do this!"})
		return
	}
	args := db.CreateMoneyRecordParams{
		CustomerID: account.CustomerID,
		Status:     "pending",
		Amount:     obj.Amount,
		Reference:  obj.Reference,
	}
	_, err = a.server.queries.CreateMoneyRecord(context.Background(), args)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record with reference exist"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	//cronJob that checks trans status

	argsBal := db.UpdateAccountBalanceManualParams{
		ID:     account.ID,
		Amount: obj.Amount,
	}
	_, err = a.server.queries.UpdateAccountBalanceManual(context.Background(), argsBal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated acct bal"})

}

type WithdrawMoneyRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Reference     string  `json:"reference" binding:"required"`
}

func (a *Account) withdrawMoney(ctx *gin.Context) {
	customerID, err := utils.GetActiveCustomer(ctx)
	if err != nil {
		return
	}

	req := new(WithdrawMoneyRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := a.server.queries.GetAccountByID(context.Background(), req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if account.CustomerID != int32(customerID) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to perform this action"})
		return
	}

	if account.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	args := db.CreateMoneyRecordParams{
		CustomerID: account.CustomerID,
		Status:     "pending",
		Amount:     -req.Amount, // Negative amount for withdrawal
		Reference:  req.Reference,
	}
	_, err = a.server.queries.CreateMoneyRecord(context.Background(), args)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record with reference exists"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	argsBal := db.UpdateAccountBalanceManualParams{
		ID:     account.ID,
		Amount: -req.Amount, // Negative amount for withdrawal
	}
	_, err = a.server.queries.UpdateAccountBalanceManual(context.Background(), argsBal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
}
