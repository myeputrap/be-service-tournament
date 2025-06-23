package usecase

import (
	"context"
	"fmt"

	"be-service-tournament/domain"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func (h *TourneyUsecase) Authorization(ctx context.Context, req domain.AuthRequest) (status int, err error) {

	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != viper.GetString("jwt.signing_method") {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(viper.GetString("jwt.signature_key")), nil
	})
	if err != nil {
		slog.Info("[Usecase][Authorization] ", "Err", err.Error())
		status = domain.StatusUnauthorized
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = domain.ErrUnauthorized
		status = domain.StatusUnauthorized
		return
	}

	// rid, ok := claims["rid"].(float64)
	// if !ok {
	// 	err = domain.ErrUnauthorized
	// 	status = domain.StatusUnauthorized
	// 	return
	// }

	loginType := claims["type"].(string)
	if loginType == "admin" {
		return
	} else if loginType != "user" {
		err = domain.ErrUnauthorized
		status = domain.StatusUnauthorized
		return
	}

	// member, err := h.EndpointRoleMember(ctx, strconv.Itoa(int(rid)), req.Method+" "+req.Path)
	// if err != nil {
	// 	slog.Info("[Usecase][Authorization] ", "Err", err.Error())
	// 	status = domain.StatusForbidden
	// 	return
	// }

	// if !member {
	// 	err = domain.ErrForbidden
	// 	status = domain.StatusForbidden
	// 	slog.Info("[Usecase][Authorization] ", "Err", err.Error())
	// 	return
	// }

	return
}
