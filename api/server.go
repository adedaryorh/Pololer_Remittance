package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golodash/galidator"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"github/adedaryorh/pooler_Remmitance_Application/utils"
	"net/http"
)

type Server struct {
	queries    *db.Store
	router     *gin.Engine
	config     *utils.Config
	tokenMaker *utils.JWTToken
}

var tokenController *utils.JWTToken
var gValid = galidator.New()

func NewServer(envPath string) *Server {
	config, err := utils.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Can not load env config: %v", err))
	}
	connection, err := sql.Open(config.DBdriver, config.DB_source+config.DB_NAME+"?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("could not load config: %v", err))
	}

	tokenController = utils.NewJWTToken(config)

	q := db.NewStore(connection)
	g := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	return &Server{
		queries: q,
		router:  g,
		config:  config,
	}
	/*
		g := gin.Default()

		g.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "welcome to the Application"})
		})

		//g.Run(fmt.Sprintf(":"))
		g.Run(fmt.Sprintf(":%v", port))
	*/

}
func (s *Server) Start(port int) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the Application"})
	})
	Customer{}.router(s)
	Authentication{}.router(s)
	Account{}.router(s)
	s.router.Run(fmt.Sprintf(":%v", port))
}
