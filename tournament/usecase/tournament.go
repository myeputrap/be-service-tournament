package usecase

import (
	"be-service-tournament/domain"
	"context"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type TourneyUsecase struct {
	mysqlRepository domain.SQLTournamentRepository
}

func NewTourneyUsecase(mysqlRepo domain.SQLTournamentRepository) domain.TournamentUsecase {
	return &TourneyUsecase{
		mysqlRepository: mysqlRepo,
	}
}

func (u *TourneyUsecase) Login(ctx context.Context, req domain.RequestLogin) (r map[string]interface{}, status int, err error) {
	slog.Info("[Usecase][Login] Login")
	user, status, err := u.mysqlRepository.GetUserLogin(ctx, req)
	if err != nil {
		slog.Error("[Usecase][Login]" + err.Error())
		return
	}

	slog.Debug("[Usecase][Login] Get user from DB success")
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		slog.Error("[Usecase][Login] Wrong username/password")
		status = domain.StatusNotFound
		return
	}
	slog.Debug("[Usecase][Login] Bcrypt compare hash and password success")
	expSecond := time.Duration(viper.GetInt64("jwt.expiration")) * time.Second
	expTime := time.Now().Add(expSecond)
	expTimeStr := expTime.Format("2006-01-02T15:04:05-0700")
	var typeRole string
	if user.RoleID == 1 {
		typeRole = "admin"
	} else {
		typeRole = "user"
	}
	claims := jwt.MapClaims{
		"id":         user.ID,
		"rid":        user.RoleID,
		"type":       typeRole,
		"token_type": "permanent",
		"exp":        expTime.Unix(),
	}
	log.Println("EXPIRED", strconv.FormatInt(expTime.Unix(), 10))
	log.Println("NOW", strconv.FormatInt(time.Now().Unix(), 10))
	// Create token
	t := jwt.NewWithClaims(jwt.GetSigningMethod(viper.GetString("jwt.signing_method")), claims)

	// Generate encoded token and send it as response.
	token, err := t.SignedString([]byte(viper.GetString("jwt.signature_key")))
	if err != nil {
		slog.Error("[Usecase][Login] Error login", "Err", err.Error())
		return
	}
	// scope, err := u.mysqlRepository.GetScopeDescriptionByRoleID(ctx, user.RoleID)
	r = map[string]interface{}{
		"token":   token,
		"expired": expTimeStr,
		//"scopes":  scope, // Add the array of strings here
		"email": user.Email,
		"name":  user.FullName,
	}
	// if err != nil {
	// 	slog.Error("[Usecase][Login] Error login", "Err", err.Error())
	// 	return
	// }
	status = domain.StatusSuccessLogin
	return
}
