package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unautorized request"})
			c.Abort()
			return
		}

		tokenSplit := strings.Split(token, "")
		if len(tokenSplit) != 2 || strings.ToLower(tokenSplit[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token, expects bearer token"})
			c.Abort()
			return
		}
		customeID, err := tokenController.VerifyToken(tokenSplit[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		c.Set("customer_id", customeID)
	}
}
