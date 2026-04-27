package server

import (
	"database/sql"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"expenses/internal/config"
	"expenses/internal/route"
)

func New(_ *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())

	route.Register(e, db)

	return e
}
