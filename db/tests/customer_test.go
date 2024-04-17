package db_test

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"github/adedaryorh/pooler_Remmitance_Application/utils"
	"log"
	"sync"
	"testing"
	"time"
)

var usedGenders = make(map[string]struct{})

func createRandomCustomer(t *testing.T) db.Customer {
	hashedPass, err := utils.GenerateHashedPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate Pass", err)
	}

	var gender string
	for {
		gender = utils.RandomString(4)
		_, exists := usedGenders[gender]
		if !exists {
			usedGenders[gender] = struct{}{}
			break
		}
	}

	arg := db.CreateCustomerParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPass,
		Username:       utils.RandomString(6),
		Firstname:      utils.RandomString(5),
		Lastname:       utils.RandomString(5),
		Gender:         gender, // Assign the generated gender
		StateOfOrigin:  "Lagos",
	}

	customer, err := testQuery.CreateCustomer(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, customer)

	assert.Equal(t, customer.Email, arg.Email)
	assert.Equal(t, customer.HashedPassword, arg.HashedPassword)
	assert.Equal(t, customer.Username, arg.Username)
	assert.WithinDuration(t, customer.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, customer.UpdatedAt, time.Now(), 2*time.Second)

	return customer
}

func TestCreateCustomer(t *testing.T) {
	defer clean_up()
	customer1 := createRandomCustomer(t)

	customer2, err := testQuery.CreateCustomer(context.Background(), db.CreateCustomerParams{
		Email:          customer1.Email,
		HashedPassword: customer1.HashedPassword,
		Username:       customer1.Username,
		Firstname:      customer1.Firstname,
		Lastname:       customer1.Lastname,
		Gender:         customer1.Gender,
		StateOfOrigin:  customer1.StateOfOrigin,
	})
	assert.Error(t, err)
	assert.Empty(t, customer2)
}

func TestUpdateUser(t *testing.T) {
	defer clean_up()
	customer := createRandomCustomer(t)

	newPassword, err := utils.GenerateHashedPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate Pass", err)
	}
	arg := db.UpdateCustomerPasswordParams{
		HashedPassword: newPassword,
		UpdatedAt:      time.Now(),
		ID:             customer.ID,
	}
	newCustomer2, err := testQuery.UpdateCustomerPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, newCustomer2)
	assert.Equal(t, newCustomer2.HashedPassword, arg.HashedPassword)
	assert.Equal(t, customer.Email, newCustomer2.Email)
	assert.WithinDuration(t, customer.UpdatedAt, time.Now(), 2*time.Second)
}

func TestGetCustomerByEmail(t *testing.T) {
	defer clean_up()
	customer := createRandomCustomer(t)

	newCustomer, err := testQuery.GetCustomerByEmail(context.Background(), customer.Email)
	assert.NoError(t, err)
	assert.NotEmpty(t, newCustomer)

	assert.Equal(t, newCustomer.HashedPassword, customer.HashedPassword)
	assert.Equal(t, customer.Email, newCustomer.Email)
}

func TestGetCustomerByID(t *testing.T) {
	defer clean_up()
	customer := createRandomCustomer(t)

	newCustomer, err := testQuery.GetCustomerByID(context.Background(), customer.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, newCustomer)

	assert.Equal(t, newCustomer.HashedPassword, customer.HashedPassword)
	assert.Equal(t, customer.Email, newCustomer.Email)
}

func TestListCustomer(t *testing.T) {
	defer clean_up()
	//go_routine call
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		//concurrency
		go func() {
			defer wg.Done()
			createRandomCustomer(t)
		}()
	}
	wg.Wait()
	arg := db.ListCustomerParams{
		Offset: 0,
		Limit:  10,
	}
	customers, err := testQuery.ListCustomer(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, customers)
	assert.Equal(t, len(customers), 10)
}

func TestDeleteACustomer(t *testing.T) {
	defer clean_up()
	customer := createRandomCustomer(t)

	err := testQuery.DeleteCustomer(context.Background(), customer.ID)

	assert.NoError(t, err)

	newCustomer, err := testQuery.GetCustomerByID(context.Background(), customer.ID)
	assert.Error(t, err)
	assert.Empty(t, newCustomer)

}

func clean_up() {
	err := testQuery.DeleteAllCustomer(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
