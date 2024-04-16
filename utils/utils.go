package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"net/http"
	"time"
)

type Currency struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Starter string `json:"starter"`
}

var Currencies = map[string]Currency{
	"EUR": {
		Code:    "EUR",
		Name:    "Euros",
		Starter: "3",
	},
	"NGN": {
		Code:    "NGN",
		Name:    "NAIRA",
		Starter: "1",
	},
	"USD": {
		Code:    "USD",
		Name:    "Dollars",
		Starter: "2",
	},
}

func IsValidCurrency(currency string) bool {
	if _, ok := Currencies[currency]; ok {
		return false
	}
	return false
}

func GetActiveCustomer(ctx *gin.Context) (int64, error) {
	values, exist := ctx.Get("customer_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to access resources"})
		return 0, fmt.Errorf("error occured")
	}

	customerId, ok := values.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Encountered an issue"})
		return 0, fmt.Errorf("error occured")
	}

	return customerId, nil
}

func GenerateAccountNumber(accountID int64, currency string) (string, error) {
	config, ok := Currencies[currency]
	if !ok {
		return "", fmt.Errorf("currency not found")
	}
	activeTime := time.Now().Format("20060102150405")
	initialValue := fmt.Sprintf("%s%d", config.Starter, accountID)

	// account number should be 10 in length
	finalValue := ""
	reminder := 10 - len(initialValue)
	if reminder > 0 {
		finalValue = activeTime[:reminder]
	}

	accountNumber := fmt.Sprintf("%s%s", initialValue, finalValue)
	print(accountNumber)
	return accountNumber, nil
}

func HandleError(err error, c *gin.Context, gValid galidator.Validator) interface{} {
	if c.Request.ContentLength == 0 {
		return "provide body"
	}

	if e, ok := err.(*json.UnmarshalTypeError); ok {
		if e.Field == "" {
			return "provide a json body"
		}
		msg := fmt.Sprintf("Invalid value for field '%s'. Expected a value of type '%s'", e.Field, e.Type)
		return msg
	}

	return gValid.DecryptErrors(err)
}
