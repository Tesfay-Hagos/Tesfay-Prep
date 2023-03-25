package api

import (
	"fmt"
	db "tesfayprep/simplebank/db/sqlc"
	"tesfayprep/simplebank/util"
	"tesfayprep/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setuproute()
	return server, nil
}

func (server *Server) setuproute() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authroutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroutes.POST("/accounts", server.createAccount)
	authroutes.GET("/accounts/:id", server.getAccount)
	authroutes.GET("/accounts", server.listAccount)
	authroutes.POST("/transfers", server.createTransfer)
	server.router = router
}
