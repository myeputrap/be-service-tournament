package http

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"errors"

	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func MiddlewareJWTAuthorizationUser(c *fiber.Ctx) error {
	slog.Info("[Validator][MiddlewareJWTAuthorization] Authorizing user/device")
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("couldn't parse claims")
		slog.Info("[Validator][MiddlewareJWTAuthorization] couldn't parse claims", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	expire, err := claims.GetExpirationTime()
	if err != nil {
		slog.Info("[Validator][MiddlewareJWTAuthorization] couldn't parse claims", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	if expire.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		slog.Info("[Validator][MiddlewareJWTAuthorization] token expired", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)

	}
	loginType := claims["type"].(string)
	if loginType != "user" {
		err = domain.ErrUnauthorized
		slog.Info("[Validator][MiddlewareJWTAuthorization] token expired", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	errSet := setPlayloadToContext(c)
	if errSet != nil {
		slog.Error("[Validator][MiddlewareJWTAuthorization] Error set playload to context", "", errSet.Error())
	}
	return c.Next()
}

func MiddlewareJWTAuthorizationDevice(c *fiber.Ctx) error {
	slog.Info("[Validator][MiddlewareJWTAuthorization] Authorizing user/user")
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("couldn't parse claims")
		slog.Info("[Validator][MiddlewareJWTAuthorization] couldn't parse claims", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	expire, err := claims.GetExpirationTime()
	if err != nil {
		slog.Info("[Validator][MiddlewareJWTAuthorization] couldn't parse claims", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	if expire.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		slog.Info("[Validator][MiddlewareJWTAuthorization] token expired", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)

	}
	loginType := claims["type"].(string)
	if loginType != "user" {
		err = domain.ErrUnauthorized
		slog.Info("[Validator][MiddlewareJWTAuthorization] token expired", "Err", err.Error())
		response := helper.NewResponse(domain.StatusUnauthorized, err.Error(), nil, nil)
		return c.Status(domain.GetHttpStatusCode(domain.StatusUnauthorized)).JSON(response)
	}
	errSet := setPlayloadToContext(c)
	if errSet != nil {
		slog.Error("[Validator][MiddlewareJWTAuthorization] Error set playload to context", "", errSet.Error())
	}

	return c.Next()
}

func setPlayloadToContext(c *fiber.Ctx) (err error) {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("couldn't parse claims")
		slog.Info("[Validator][MiddlewareJWTAuthorization] couldn't parse claims", "Err", err.Error())
	}
	var userData helper.UserLogin
	id := claims["id"].(float64)
	userData.ID = int32(id)
	if claims["rid"] != nil {
		rid := claims["rid"].(float64)
		ridInt := int32(rid)
		userData.RID = &ridInt
	}
	if claims["hid"] != nil {
		hid := claims["hid"].(float64)
		hidInt := int32(hid)
		userData.HID = &hidInt
	}
	c.Locals("user_login", userData)
	return
}
