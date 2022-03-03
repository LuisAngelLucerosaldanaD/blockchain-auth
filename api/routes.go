package api

import (
	"blion-auth/api/handlers/login"
	"github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// @title bjungle Blockchain API.
// @version 1.0
// @description Documentation Api Blockchain.
// @termsOfService https://www.bjungle-id.com/terms/
// @contact.name API Support.
// @contact.email info@bjungle-id.com
// @license.name Software Owner
// @license.url https://www.bjungle-id.com/licenses/LICENSE-1.0.html
// @host localhost:50025
// @BasePath /
func routes(db *sqlx.DB, loggerHttp bool, allowedOrigins string) *fiber.App {
	app := fiber.New()

	prometheus := fiberprometheus.New("nexumApiAuth")
	prometheus.RegisterAt(app, "/metrics")

	//app.Get("/swagger/*", fiberSwagger.Handler)
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	app.Use(recover.New())
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST",
	}))
	if loggerHttp {
		app.Use(logger.New())
	}
	TxID := uuid.New().String()

	loadRoutes(app, db, TxID)

	return app
}

func loadRoutes(app *fiber.App, db *sqlx.DB, TxID string) {
	/*users.RouterCreateUser(app, db, TxID)
	wallets.RouterWallet(app, db, TxID)
	transaction.RouterTransaction(app, db, TxID)
	blocks.RouterBlock(app, db, TxID)
	mine.RouterMine(app, db, TxID)*/
	login.RouterLogin(app, db, TxID)
	/*genesis.RouterGenesis(app, db, TxID)
	cipher.RouterCipher(app, db, TxID)
	category.RouterCategory(app, db, TxID)
	dictionary.RouterDictionary(app, db, TxID)*/
}
