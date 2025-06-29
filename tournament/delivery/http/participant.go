package http

import (
	"be-service-tournament/constant"
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
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

func (h *tournamentHandler) UpdateParticipantStatus(c *fiber.Ctx) error {
	slog.Info("[Handler][UpdateParticipantStatus] UpdateParticipantStatus")

	var status int
	var userID int64
	var err error
	var message string
	var req domain.UpdateParticipantStatusRequest
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
			slog.Error("[Handler][UpdateParticipantStatus] " + message + ": " + err.Error())
			return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
		}
	} else {
		status = domain.StatusMissingParameter
		message = domain.GetCustomStatusMessage(status, "id")
		slog.Error("[Handler][UpdateParticipantStatus] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	req.UserID = userID
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][UpdateParticipantStatus]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	if !helper.ContainsAnySingle(req.Status, registrationStatus) {
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Value of status is not recognized", nil, nil))
	}
	status, err = h.hospitalityusecase.UpdateParticipant(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][UpdateParticipantStatus] Error UpdateParticipantStatus", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't UpdateParticipantStatus"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, nil)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) PostImagePaymentProof(c *fiber.Ctx) error {
	slog.Info("[Handler][PostImagePaymentProof] PostImagePaymentProof")

	var status int
	var parID int64
	var err error
	var message string
	var req domain.RequestPaymentProffImage
	idStr := c.Params("id")
	if idStr != "" {
		parID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			status = domain.StatusWrongValue
			message = domain.GetCustomStatusMessage(status, "id")
			slog.Error("[Handler][PostImagePaymentProof] " + message + ": " + err.Error())
			return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
		}
	} else {
		status = domain.StatusMissingParameter
		message = domain.GetCustomStatusMessage(status, "id")
		slog.Error("[Handler][PostImagePaymentProof] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	req.ParticipantID = parID
	file, err := c.FormFile("images")
	if err != nil {
		slog.Error("[Handler][PostImagePaymentProof] " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get uploaded image",
		})
	}
	req.Images = file
	if req.Images.Size > constant.MaxSizeHotelBackground {
		slog.Error("[Handler][CreateTheme] The image must be less than 15Mb", "Err", nil)
		status = domain.StatusWrongValue
		message = "file background must be less than 100Mb"
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}

	err = helper.CheckFile(req.Images)
	if err != nil && errors.Is(err, domain.ErrBadRequest) {
		slog.Error("[Handler][CreateTheme] check file", "Err", err.Error())
		status = domain.StatusWrongValue
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))

	} else if err != nil {
		slog.Error("[Handler][CreateTheme] check file", "Err", err.Error())
		status = domain.StatusInternalServerError
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Internal server error", nil, nil))
	}
	status, err = h.hospitalityusecase.CreatePaymentProofImage(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][UpdateParticipant] Error UpdateParticipant", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't UpdateParticipant"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, nil)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) GetAllPaticipant(c *fiber.Ctx) error {
	var err error
	var req domain.GetAllParticipantRequest
	if err := c.QueryParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][GetAllPaticipant]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	//checkValue
	req.Limit, err = strconv.ParseInt(c.Query("limit", strconv.FormatInt(int64(viper.GetInt("default_limit_query")), 10)), 10, 64)
	if err != nil {
		return c.Status(status).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}

	req.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return c.Status(status).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if req.Sort == "" {
		req.Sort = "id"
	}
	if req.Order == "" {
		req.Order = "asc"
	}
	r, status, err := h.hospitalityusecase.GetAllParticipant(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][Login] Error GetAllPaticipant", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't GetAllPaticipant"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}
