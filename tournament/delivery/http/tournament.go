package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"log/slog"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type tournamentHandler struct {
	hospitalityusecase domain.TournamentUsecase
	validator          *validator.Validate
}

func NewHandler(u domain.TournamentUsecase, validator *validator.Validate) *tournamentHandler {
	return &tournamentHandler{u, validator}
}

func (h *tournamentHandler) Login(c *fiber.Ctx) error {
	var status int
	var message string
	var req domain.RequestLogin
	if err := c.BodyParser(&req); err != nil {
		err = domain.ErrBadRequest
		slog.Error("[Handler][Login] Error login", "Err", err.Error())
		status = domain.StatusInternalServerError
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if req.Email == "" || req.Password == "" {
		status = domain.StatusMissingParameter
		message = domain.GetCustomStatusMessage(status, "email/password")
		slog.Error("[Handler][Login] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		status = domain.StatusInvalidEmailPassword
		message = domain.GetCustomStatusMessage(status, "")
		slog.Error("[Handler][Login] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	r, status, err := h.hospitalityusecase.Login(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][Login] Error login", "Err", err.Error())
		switch status {
		case domain.StatusNotFound:
			status = domain.StatusInvalidEmailPassword
			message = domain.GetCustomStatusMessage(status, "")
		case domain.StatusUnauthorizedUnverified:
			status = domain.StatusUnauthorizedUnverified
			message = domain.GetCustomStatusMessage(status, "")
		case domain.StatusUnauthorizedBlockedAccount:
			status = domain.StatusUnauthorizedBlockedAccount
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) TestAdmin(c *fiber.Ctx) error {
	response := helper.NewResponse(domain.StatusSuccess, "OK", nil, nil)
	return c.Status(domain.GetHttpStatusCode(domain.StatusSuccess)).JSON(response)
}

func (h *tournamentHandler) InquiryTourneyPublic(c *fiber.Ctx) error {
	var status int
	var message string
	r, status, err := h.hospitalityusecase.InquiryTourneyPublic(c.Context())
	if err != nil {
		slog.Error("[Handler][InquiryTourneyPublic] Error InquiryTourneyPublic", "Err", err.Error())
		switch status {
		case domain.StatusNotFound:
			status = domain.StatusNotFound
			message = domain.GetCustomStatusMessage(status, "")
		default:
			status = domain.StatusInternalServerError
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}
