package main

import (
	"net/http"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/api"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/config"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/connection"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/repository"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/service"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	// Setup JWT Middleware
	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Return 401 Unauthorized if JWT validation fails
			return ctx.Status(http.StatusUnauthorized).
				JSON(dto.CreateResponseError("Unauthorized"))
		},
	})

	// Initialize Repositories and Services
	customerRepository := repository.NewCustomer(dbConnection)
	userRepository := repository.NewUser(dbConnection)
	bookRepository := repository.NewBook(dbConnection)
	bookStockRepository := repository.NewBookStock(dbConnection)
	journalRepository := repository.NewJournal(dbConnection)
	mediaRepository := repository.NewMedia(dbConnection)
	chargeRepository := repository.NewCharge(dbConnection)

	customerService := service.NewCustomer(customerRepository)
	authService := service.NewAuth(cnf, userRepository)
	bookService := service.NewBook(bookRepository, bookStockRepository, mediaRepository, cnf)
	bookStockService := service.NewBookStock(bookRepository, bookStockRepository)
	journalService := service.NewJournal(
		journalRepository,
		bookRepository,
		bookStockRepository,
		customerRepository,
		chargeRepository,
	)
	mediaService := service.NewMedia(cnf, mediaRepository)

	// Setup API Routes
	api.NewCustomer(app, customerService, jwtMidd)
	api.NewAuth(app, authService)
	api.NewBook(app, bookService, jwtMidd)
	api.NewBookStock(app, bookStockService, jwtMidd)
	api.NewJournal(app, journalService, jwtMidd)
	api.NewMedia(app, cnf, mediaService, jwtMidd)

	app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
