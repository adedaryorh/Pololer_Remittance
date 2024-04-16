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

type Authentication struct {
	server *Server
}

func (a Authentication) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/authentication")
	serverGroup.POST("login", a.login)
	serverGroup.POST("register", a.register)
}

func (a *Authentication) register(c *gin.Context) {
	//creating a user instance
	customer := new(CustomerParams)
	//or var user UserParams
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := utils.GenerateHashedPassword(customer.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateCustomerParams{
		Email:          customer.Email,
		HashedPassword: hashedPassword,
		Username:       customer.Username,
	}

	newUser, err := a.server.queries.CreateCustomer(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, CustomerResponse{}.toCustomerResponse(&newUser))
}

func (a Authentication) login(c *gin.Context) {
	customer := new(CustomerParams)
	eViewer := gValid.Validator(CustomerParams{})

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": eViewer.DecryptErrors(err)})
		return
	}

	dbCustomer, err := a.server.queries.GetCustomerByEmail(context.Background(), customer.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or pass"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.VerifyPassword(customer.HashedPassword, dbCustomer.HashedPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or pass"})
		return
	}
	token, err := tokenController.CreateToken(dbCustomer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
