package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"be-service-tournament/tournament/delivery/middleware/authorization"
	"log/slog"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func RouterAPI(app *fiber.App, us domain.TournamentUsecase) {
	requestValidator := validator.New()
	basePath := viper.GetString("server.http.base_path")

	handler := &tournamentHandler{us, requestValidator}
	authMiddleware := authorization.New(authorization.Config{Usecase: us})
	api := app.Group(basePath)
	adm := app.Group("/adm")
	user := app.Group("/user")
	adm.Post("/login", handler.Login)
	api.Get("/tourney", handler.InquiryTourneyPublic)

	adm.Use(jwtware.New(JWTMiddlewareConfiguration()), MiddlewareJWTAuthorizationAdmin)
	user.Use(jwtware.New(JWTMiddlewareConfiguration()), MiddlewareJWTAuthorizationUser)
	adm.Get("/test-admin", authMiddleware, handler.TestAdmin)
	//user.Get("/test-user", authMiddleware, handler.TestUser)
}

func JWTMiddlewareConfiguration() jwtware.Config {
	slog.Debug("[Router][JWTMiddlewareConfiguration] Build configuration for JWT middleware")
	var config jwtware.Config
	config.SigningKey = jwtware.SigningKey{Key: []byte(viper.GetString("jwt.signature_key"))}
	config.ErrorHandler = JWTErrorHandler
	config.Filter = JWTFilter
	return config
}

func JWTErrorHandler(c *fiber.Ctx, err error) error {
	slog.Error("[Router][JWTErrorHandler] Unauthorized", "Err", "Invalid token")
	response := helper.NewResponse(domain.StatusUnauthorized, "Invalid token", nil, nil)
	return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
}

func JWTFilter(c *fiber.Ctx) bool {
	return c.Path() == "/login"
}
