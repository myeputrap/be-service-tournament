package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var registrationStatus = []string{
	"Applied",
	"Shortlisted",
	"Cek Payment",
	"Registered",
	"Waiting list",
	"Denied",
}

func (h *tournamentHandler) CreateUserParticipant(c *fiber.Ctx) error {
	slog.Info("[Handler][CreateUserParticipant] CreateUserParticipant")

	var status int
	var message string
	var req domain.ParticipantDTO

	if err := c.BodyParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][CreateUserParticipant]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	status, err := h.hospitalityusecase.FormPartnershipParticipant(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][CreateUserParticipant] Error CreateUserParticipant", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't CreateUserParticipant"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccessCreate
	response := helper.NewResponse(status, "OK", nil, nil)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) UpdateParticipant(c *fiber.Ctx) error {
	slog.Info("[Handler][UpdateParticipant] UpdateParticipant")

	var status int
	var userID int64
	var err error
	var message string
	var req domain.UpdateParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	idStr := c.Params("id")
	if idStr != "" {
		userID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			status = domain.StatusWrongValue
			message = domain.GetCustomStatusMessage(status, "id")
			slog.Error("[Handler][UpdateParticipant] " + message + ": " + err.Error())
			return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
		}
	} else {
		status = domain.StatusMissingParameter
		message = domain.GetCustomStatusMessage(status, "id")
		slog.Error("[Handler][UpdateParticipant] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	req.UserID = userID
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][UpdateParticipant]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	if !helper.ContainsAnySingle(req.Status, registrationStatus) {
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of status is not recognized", nil, nil))
	}
	status, err = h.hospitalityusecase.UpdateParticipant(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][UpdateParticipant] Error UpdateParticipant", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't UpdateParticipant"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccessCreate
	response := helper.NewResponse(status, "OK", nil, nil)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}
