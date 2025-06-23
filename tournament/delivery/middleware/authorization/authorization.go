package authorization

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Usecase domain.TournamentUsecase
}

var ConfigDefault = Config{
	Usecase: nil,
}

func configDefault(config Config) {
	ConfigDefault.Usecase = config.Usecase
}

func New(config Config) fiber.Handler {
	configDefault(config)
	return AuthorizationMiddleware
}

func AuthorizationMiddleware(c *fiber.Ctx) (err error) {
	slog.Info("[Middleware][Authorization]")
	var status int
	userLogin := helper.GetUserLogin(c.Context())
	var roleID string
	if userLogin.RID != nil && *userLogin.RID > 0 {
		roleID = strconv.FormatInt(int64(*userLogin.RID), 10)
	} else {
		err = errors.New("no role id found")
		slog.Info("[Authorization][AuthorizationMiddleware] ", "Err", err.Error())
		status = domain.StatusForbidden
		response := helper.NewResponse(status, domain.GetCustomStatusMessage(status, ""), nil, nil)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
	}
	path := c.Route().Path
	method := c.Route().Method
	val := method + " " + path
	slog.Info("VALUE: ", "", val)
	slog.Info("VALUE r: ", "", roleID)
	// member, err := ConfigDefault.Usecase.EndpointRoleMember(c.Context(), roleID, val)
	// if err != nil {
	// 	slog.Info("[Authorization][AuthorizationMiddleware] ", "Err", err.Error())
	// 	status = domain.StatusForbidden
	// 	response := helper.NewResponse(status, domain.GetCustomStatusMessage(status, ""), nil, nil)
	// 	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
	// }
	// if !member {
	// 	err = errors.New("forbidden")
	// 	slog.Info("[Authorization][AuthorizationMiddleware] ", "Err", err.Error())
	// 	status = domain.StatusForbidden
	// 	response := helper.NewResponse(status, domain.GetCustomStatusMessage(status, ""), nil, nil)
	// 	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
	// }
	return c.Next()
}
