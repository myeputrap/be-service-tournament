package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"log/slog"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var genderEligibility = []string{
	"Male only",
	"Female only",
	"Mixed",
	"Free",
}

var typeTournament = []string{
	"Elimination",
	"Elimination&Group",
}
var status int
var message string

type tournamentHandler struct {
	hospitalityusecase domain.TournamentUsecase
	validator          *validator.Validate
}

func NewHandler(u domain.TournamentUsecase, validator *validator.Validate) *tournamentHandler {
	return &tournamentHandler{u, validator}
}

func (h *tournamentHandler) Login(c *fiber.Ctx) error {
	slog.Info("[Handler][Login] Login")
	var req domain.RequestLogin
	var message string
	if err := c.BodyParser(&req); err != nil {
		err = domain.ErrBadRequest
		slog.Error("[Handler][Login] Error login", "Err", err.Error())
		status = domain.StatusInternalServerError
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][Login]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
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
		if status == domain.StatusInternalServerError {
			message = "something wrong with login"
		} else {
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
	slog.Info("[Handler][InquiryTourneyPublic] InquiryTourneyPublic")
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

func (h *tournamentHandler) CreateTournament(c *fiber.Ctx) error {
	slog.Info("[Handler][CreateTournament] CreateTournament")
	var req domain.Tournament
	if err := c.BodyParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][Login]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	//check value
	if !helper.ContainsAnySingle(req.GenderEligibility, genderEligibility) {
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of gender GenderEligibility is not recognized", nil, nil))
	}
	if !helper.ContainsAnySingle(req.Type, typeTournament) {
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of gender TournamentType is not recognized", nil, nil))
	}

	if req.Type == "Elimination" || (req.Quota != 16 && req.Quota != 32 && req.Quota != 64) {
		status = domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of quota is not recognized", nil, nil))
	}

	if req.Type == "GroupElimination&Group" || (req.Quota != 32 && req.Quota != 64) {
		status = domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of quota is not recognized for tournament type", nil, nil))
	}
	userData := helper.GetUserLogin(c.Context())
	req.CreatedBy = int64(userData.ID)
	r, status, err := h.hospitalityusecase.CreateTournament(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][Login] Error CreateUser", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't CreateUser"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccessCreate
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}
