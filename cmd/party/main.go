package main

import (
	"log"

	"github.com/alesr/getground/internal/app"
	"github.com/alesr/getground/internal/app/partyctrl"
	"github.com/alesr/getground/internal/pkg/party"
	"github.com/alesr/getground/internal/pkg/party/repository"
	"github.com/alesr/getground/pkg/database"
	fiber "github.com/gofiber/fiber/v2"
	envars "github.com/netflix/go-env"
	"go.uber.org/zap"
)

type config struct {
	Env            string `env:"ENV,default=dev"`
	AppName        string `env:"APP_NAME,default=getground"`
	Port           string `env:"PORT,default=3000"`
	AllowHeaders   string `env:"ALLOW_HEADERS,default=Origin, Content-Type, Accept"`
	PartyTableSize int    `env:"PARTY_TABLE_SIZE,default=12"`
	DB             dbConfig
}

type dbConfig struct {
	Host     string `env:"MYSQL_HOST,default=127.0.0.1"`
	User     string `env:"MYSQL_USER,default=user"`
	Password string `env:"MYSQL_PASSWORD,default=password"`
	Name     string `env:"MYSQL_DB,default=party_db"`
}

func newConfig() *config {
	var cfg config
	if _, err := envars.UnmarshalFromEnviron(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func main() {
	// Load configuration
	cfg := newConfig()

	// Initalize logger

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to initialize production logger:", err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	logger.Named(cfg.AppName)

	logger.Info("initializing app", zap.String("env", cfg.Env))

	dbConn, err := database.Connection(cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		logger.Fatal("failed to get database connection", zap.Error(err))
	}

	logger.Info("applying database migrations")

	dbConn.Table("guests").AutoMigrate(&repository.Guest{})
	dbConn.Table("tables").AutoMigrate(&repository.Table{})

	logger.Info("database migrations DONE")

	// Initialize repository
	partyRepo := repository.New(logger, dbConn)

	// Initialize service
	partyService := party.New(logger, partyRepo, cfg.PartyTableSize)

	// Initialize HTTP router

	fiberApp := fiber.New()

	restApp := app.New(
		logger,
		fiberApp,
		partyctrl.New(logger, partyService),
	)

	if err := restApp.Run(cfg.Port); err != nil {
		logger.Fatal("failed to run rest app", zap.Error(err))
	}
}
