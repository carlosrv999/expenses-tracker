package server

import (
	"database/sql"
	"net/http" // ← NEW IMPORT

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"expenses/internal/config"
	"expenses/internal/route"
)

func New(_ *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()

	// ✅ CORS for your localhost frontend
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // React / Create-React-App
			"http://127.0.0.1:3000",
			"http://localhost:5173", // Vite (most common now)
			"http://127.0.0.1:5173",
			// Add any other frontend ports you use
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true, // keep this if you use cookies / auth
	}))

	e.Use(middleware.Recover())

	route.Register(e, db)

	return e
}
