package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *tournamentHandler) CreateUser(c *fiber.Ctx) error {
	slog.Info("[Handler][CreateUser] CreateUser")

	var status int
	var message string
	var req domain.UserRequestDTO

	if err := c.BodyParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][CreateUser]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	r, multiErr, status, err := h.hospitalityusecase.CreateUser(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][CreateUser] Error CreateUser", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't CreateUser"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, multiErr))
	}
	status = domain.StatusSuccessCreate
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) ListUserPartner(c *fiber.Ctx) error {
	slog.Info("[Handler][ListUserPartner] ListUserPartner")
	var status int
	var message string
	var req domain.GetAllUserRequestPartner
	if err := c.QueryParser(&req); err != nil {
		status = domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}
	userData := helper.GetUserLogin(c.Context())
	req.UIDSearcher = int64(userData.ID)
	r, status, err := h.hospitalityusecase.GetUserPartner(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][Login] Error ListUserPartner", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't ListUserPartner"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) CreateAdmin(c *fiber.Ctx) error {
	slog.Info("[Handler][CreateUser] CreateUser")

	var status int
	var message string
	var req domain.UserRequestDTO

	if err := c.BodyParser(&req); err != nil {
		status := domain.StatusBadRequest
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, err.Error(), nil, nil))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.Error("[Handler][LoCreateAdmingin]", "Err", err)
		status = domain.StatusMissingParameter
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, "Bad Request", nil, nil))
	}
	r, multiErr, status, err := h.hospitalityusecase.CreateAdmin(c.Context(), req)
	if err != nil {
		slog.Error("[Handler][CreateAdmin] Error CreateUser", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't CreateUser"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, multiErr))
	}
	status = domain.StatusSuccessCreate
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) GetUserByID(c *fiber.Ctx) error {
	slog.Info("[Handler][GetUserByID] GetUserByID")
	var status int
	var userID int64
	var err error
	var message string
	idStr := c.Params("id")
	if idStr != "" {
		userID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			status = domain.StatusWrongValue
			message = domain.GetCustomStatusMessage(status, "id")
			slog.Error("[Handler][GetUserByID] " + message + ": " + err.Error())
			return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
		}
	} else {
		status = domain.StatusMissingParameter
		message = domain.GetCustomStatusMessage(status, "id")
		slog.Error("[Handler][GetUserByID] " + message)
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	r, status, err := h.hospitalityusecase.GetUserByID(c.Context(), userID)
	if err != nil {
		slog.Error("[Handler][CreateAdmin] Error GetUserByID", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't GetUserByID"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) GetUserByDetail(c *fiber.Ctx) error {
	slog.Info("[Handler][GetUserByDetail] GetUserByDetail")
	var status int
	var userID int64
	var err error
	var message string
	dataLogin := helper.GetUserLogin(c.Context())
	userID = int64(dataLogin.ID)
	r, status, err := h.hospitalityusecase.GetUserByID(context.Background(), userID)
	if err != nil {
		slog.Error("[Handler][CreateAdmin] Error GetUserByDetail", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't GetUserByDetail"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}

func (h *tournamentHandler) CreateImagePaymentProff(c *fiber.Ctx) error {
	slog.Info("[Handler][CreateImagePaymentProff] CreateImagePaymentProff")
	var status int
	var userID int64
	var err error
	var message string
	dataLogin := helper.GetUserLogin(c.Context())
	userID = int64(dataLogin.ID)
	r, status, err := h.hospitalityusecase.GetUserByID(context.Background(), userID)
	if err != nil {
		slog.Error("[Handler][CreateAdmin] Error CreateImagePaymentProff", "Err", err.Error())
		if status == domain.StatusInternalServerError {
			message = "something wrong, can't CreateImagePaymentProff"
		} else {
			message = domain.GetCustomStatusMessage(status, "")
		}
		return c.Status(domain.GetHttpStatusCode(status)).JSON(helper.NewResponse(status, message, nil, nil))
	}
	status = domain.StatusSuccess
	response := helper.NewResponse(status, "OK", nil, r)
	return c.Status(domain.GetHttpStatusCode(status)).JSON(response)
}
