package api

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"github/adedaryorh/pooler_Remmitance_Application/utils"
	"net/http"
	"time"
)

type Customer struct {
	server *Server
}

func (u Customer) router(server *Server) {
	u.server = server
	//AuthenticatedMiddleware()
	serverGroup := server.router.Group("/customer")
	serverGroup.GET("", u.listCustomers)
	serverGroup.GET("me", u.getLoggedInCustomer)
	serverGroup.POST("createCustomer", u.createCustomer)
}

type CustomerParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Gender         string `json:"gender"`
	StateOfOrigin  string `json:"state_of_origin"`
}

func (u *Customer) listCustomers(c *gin.Context) {
	arg := db.ListCustomerParams{
		Offset: 0,
		Limit:  10,
	}
	customers, err := u.server.queries.ListCustomer(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newcustomer := []CustomerResponse{}

	for _, k := range customers {
		o := CustomerResponse{}.toCustomerResponse(&k)
		newcustomer = append(newcustomer, *o)
	}

	c.JSON(http.StatusOK, newcustomer)
}

func (u *Customer) createCustomer(c *gin.Context) {
	var customer CustomerParams
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(customer.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	argument := db.CreateCustomerParams{
		Email:          customer.Email,
		HashedPassword: hashedPassword,
		Username:       customer.Username,
		Firstname:      customer.Firstname,
		Lastname:       customer.Lastname,
		Gender:         customer.Gender,
		StateOfOrigin:  customer.Gender,
	}
	newcustomer, err := u.server.queries.CreateCustomer(context.Background(), argument)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "customer exit in system"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, CustomerResponse{}.toCustomerResponse(&newcustomer))
}

func (u Customer) getLoggedInCustomer(c *gin.Context) {
	values, exist := c.Get("customer_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to access resources"})
	}

	customerid, ok := values.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Encountered an issue"})
	}

	customer, err := u.server.queries.GetCustomerByID(context.Background(), customerid)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to access resources "})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, CustomerResponse{}.toCustomerResponse(&customer))
}

type CustomerResponse struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Gender        string    `json:"gender"`
	StateOfOrigin string    `json:"state_of_origin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u CustomerResponse) toCustomerResponse(cutomer *db.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:            cutomer.ID,
		Email:         cutomer.Email,
		Gender:        cutomer.Gender,
		StateOfOrigin: cutomer.StateOfOrigin,
		CreatedAt:     cutomer.CreatedAt,
		UpdatedAt:     cutomer.UpdatedAt,
	}
}
