package main

import (
	"database/sql"

	"github.com/isAshithKumarGowda/Echo_Auth/internals/handlers"
	"github.com/isAshithKumarGowda/Echo_Auth/repository"
	"github.com/labstack/echo/v4"
)

func InitialiseHttpRouter(db *sql.DB) *echo.Echo {
	e := echo.New()
	authHandler := handlers.NewAuthHandler(repository.NewAuthRepo(db))
	e.POST("/:type/register", authHandler.RegisterHandler)
	e.POST("/:type/login", authHandler.LoginHandler)
	e.POST("/:type/verifyEmail", authHandler.VeirfyEmail)
	return e
}
