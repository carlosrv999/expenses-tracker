package route

import (
	"database/sql"

	"github.com/labstack/echo/v5"

	"expenses/internal/handler"
	"expenses/internal/handler/category"
	"expenses/internal/handler/expense"
	paymentmethod "expenses/internal/handler/payment_method"
	"expenses/internal/handler/tag"
	"expenses/internal/repository"
	"expenses/internal/service"
)

func Register(e *echo.Echo, db *sql.DB) {
	categoryRepo := repository.NewCategoryRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	tagRepo := repository.NewTagRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	categorySvc := service.NewCategoryService(categoryRepo)
	paymentMethodSvc := service.NewPaymentMethodService(paymentMethodRepo)
	tagSvc := service.NewTagService(tagRepo)
	expenseSvc := service.NewExpenseService(expenseRepo, tagRepo)

	healthHandler := handler.NewHealthHandler(db)
	categoryHandler := category.NewHandler(categorySvc)
	paymentMethodHandler := paymentmethod.NewHandler(paymentMethodSvc)
	tagHandler := tag.NewHandler(tagSvc)
	expenseHandler := expense.NewHandler(expenseSvc)

	e.GET("/health", healthHandler.Check)

	api := e.Group("/api/v1")
	categoryHandler.Register(api.Group("/categories"))
	paymentMethodHandler.Register(api.Group("/payment-methods"))
	tagHandler.Register(api.Group("/tags"))
	expenseHandler.Register(api.Group("/expenses"))
}
