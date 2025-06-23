package main

import (
	"be-service-tournament/config"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_DeliveryHTTP "be-service-tournament/tournament/delivery/http"
	_RepoMySQL "be-service-tournament/tournament/repository/mysql"
	_Usecase "be-service-tournament/tournament/usecase"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	//_DeliveryHTTP "be-service-tournament/tournament/delivery/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func main() {
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	logLevel := &slog.LevelVar{}
	logOpt := &slog.HandlerOptions{
		Level: logLevel,
	}
	textLog := slog.NewTextHandler(os.Stdout, logOpt)
	jsonLog := slog.NewJSONHandler(os.Stdout, logOpt)
	slog.SetDefault(slog.New(textLog))

	// Config file
	config.ReadConfig(*configFile)

	// Set log type and log level
	if viper.GetString("server.log_type") == "json" {
		slog.SetDefault(slog.New(jsonLog))
	}
	switch viper.GetString("server.log_level") {
	case "error":
		logLevel.Set(slog.LevelError)
	case "warning":
		logLevel.Set(slog.LevelWarn)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	}

	if logLevel.Level().Level() == slog.LevelDebug {
		c := viper.AllSettings()
		bs, err := yaml.Marshal(c)
		if err != nil {
			slog.Error("Unable to marshal config to YAML", "error", err)
		}
		slog.Debug(string(bs))
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok search_path=%s",
		viper.GetString("postgres.host"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.database"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.schema"),
	)
	fmt.Println("==============================================================")
	fmt.Println(dsn)
	// Set up logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open DB with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("error opening DB:", err)
	}

	// Set connection pool config
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get DB from GORM:", err)
	}
	sqlDB.SetMaxOpenConns(viper.GetInt("postgres.set_max_open_conns"))
	sqlDB.SetMaxIdleConns(viper.GetInt("postgres.set_max_idle_conns"))
	sqlDB.SetConnMaxIdleTime(viper.GetDuration("postgres.set_conn_max_idle_time"))
	sqlDB.SetConnMaxLifetime(viper.GetDuration("postgres.set_conn_max_lifetime"))

	// Optional: test ping
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("failed to close DB:", err)
		}
	}()

	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.http.prefork"),
		StrictRouting: viper.GetBool("server.http.strict_routing"),
		CaseSensitive: viper.GetBool("server.http.case_sensitive"),
		BodyLimit:     viper.GetInt("server.http.body_limit"),
	})
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
	app.Get("/healthy", func(c *fiber.Ctx) error {
		return c.SendString("i'm alive")
	})
	webLogger := slog.New(textLog)
	if viper.GetString("server.log_type") == "json" {
		webLogger = slog.New(jsonLog)
	}
	app.Use(slogfiber.New(webLogger))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("server.http.allows_origin"),
	}))
	app.Use(func(c *fiber.Ctx) error {
		server := "hs/und"
		version := os.Getenv("SERVER_VERSION")
		if version != "" {
			server = "hs/" + version
		}
		// Set version on server header
		slog.Debug(server)
		c.Set("Server", server)
		slog.Debug("respheader", "respheader", c.GetRespHeaders())
		return c.Next()
	})

	// HTTP routing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World")
	})
	//register uc and repository
	mysqlRepo := _RepoMySQL.NewSQLTournamentRepository(db)
	tourneyUsecase := _Usecase.NewTournamentUsecase(mysqlRepo)
	_DeliveryHTTP.RouterAPI(app, tourneyUsecase)
	// Initialize HTTP web framework
	log.Println("HTTP server is running...")
	go func() {
		if err := app.Listen(":" + viper.GetString("server.http.port")); err != nil {
			slog.Error("Failed to listen", "port", viper.GetString("server.http.port"))
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Gracefully shutdown")
	err = app.Shutdown()
	if err != nil {
		slog.Warn("Unfortunately the shutdown wasn't smooth", "err", err)
	}
}
