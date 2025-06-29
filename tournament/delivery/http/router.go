package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"be-service-tournament/tournament/delivery/middleware/authorization"
	"fmt"
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
	api.Post("/login", handler.Login)
	api.Get("/tourney", handler.InquiryTourneyPublic)
	api.Post("/user", handler.CreateUser)
	api.Post("/admin", handler.CreateAdmin)

	adm.Use(jwtware.New(JWTMiddlewareConfiguration()), MiddlewareJWTAuthorizationAdmin)
	user.Use(jwtware.New(JWTMiddlewareConfiguration()), MiddlewareJWTAuthorizationUser)

	adm.Post("/tournament", authMiddleware, handler.CreateTournament)
	adm.Get("/tournament", authMiddleware, handler.GetAllTournament)
	adm.Get("/tournament/:id", authMiddleware, handler.GetTournamentByID)
	adm.Put("/tournament", authMiddleware, handler.UpdateTournament)
	adm.Delete("/tournament/:id", authMiddleware, handler.DeleteTournamentByID)
	adm.Get("/tournament-user/:id", authMiddleware, handler.GetTournamentParticipant)

	user.Get("/list-user-partner", authMiddleware, handler.ListUserPartner)
	user.Post("/user-participant", authMiddleware, handler.CreateUserParticipant)
	adm.Get("/user/:id", authMiddleware, handler.GetUserByID)
	user.Get("/user-tournament", authMiddleware, handler.GetUserTournament)
	user.Get("/user-detail", authMiddleware, handler.GetUserByDetail)

	user.Post("/payment-proof/:id", authMiddleware, handler.PostImagePaymentProof)

	adm.Patch("/participant/:id/status", authMiddleware, handler.UpdateParticipantStatus)
	adm.Get("/participant", authMiddleware, handler.GetAllPaticipant)

}

func JWTMiddlewareConfiguration() jwtware.Config {
	var config jwtware.Config
	config.SigningKey = jwtware.SigningKey{Key: []byte(viper.GetString("jwt.signature_key"))}
	fmt.Println("JWT secret:", viper.GetString("jwt.signature_key"))
	config.ErrorHandler = JWTErrorHandler
	config.Filter = JWTFilter
	return config
}

func JWTErrorHandler(c *fiber.Ctx, err error) error {
	slog.Error("[Router][JWTErrorHandler] Unauthorized", "Err", "Invalid token")
	response := helper.NewResponse(domain.StatusUnauthorized, "Invalid token1111111", nil, nil)
	return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
}

func JWTFilter(c *fiber.Ctx) bool {
	return c.Path() == "/login"
}
