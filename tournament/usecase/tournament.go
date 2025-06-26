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

func NewTournamentUsecase(mysqlRepo domain.SQLTournamentRepository) domain.TournamentUsecase {
	return &TourneyUsecase{
		mysqlRepository: mysqlRepo,
	}
}

func (h *TourneyUsecase) Login(ctx context.Context, req domain.RequestLogin) (r map[string]interface{}, status int, err error) {
	slog.Info("[Usecase][Login] Login")
	user, status, err := h.mysqlRepository.GetUserLogin(ctx, req)
	if err != nil {
		slog.Error("[Usecase][Login]" + err.Error())
		return
	}

	slog.Debug("[Usecase][Login] Get user from DB success")
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		slog.Error("[Usecase][Login] Wrong username/password")
		status = domain.StatusInvalidEmailPassword
		err = domain.ErrBadRequest
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

func (h *TourneyUsecase) InquiryTourneyPublic(ctx context.Context) (response []domain.InquiryTourneyPublicResponse, status int, err error) {
	slog.Info("[Usecase][InquiryTourneyPublic] InquiryTourneyPublic")
	var countParticipant int64
	tourneys, status, err := h.mysqlRepository.GetTournament(ctx)
	if err != nil {
		slog.Error("[Usecase][InquiryTourneyPublic] " + err.Error())
		status = domain.StatusNotFound
		return
	}
	response = make([]domain.InquiryTourneyPublicResponse, len(tourneys))
	for i, v := range tourneys {
		countParticipant, status, err = h.mysqlRepository.CountParticipantTournament(ctx, int32(v.ID))
		if err != nil {
			slog.Error("[Usecase][InquiryTourneyPublic] " + err.Error())
			status = domain.StatusNotFound
			return
		}
		response[i] = domain.InquiryTourneyPublicResponse{
			Name:          v.Name,
			StartDate:     v.StartDate.Format("2 Jan 2006, 15.04"),
			EndDate:       v.EndDate.Format("2 Jan 2006, 15.04"),
			Location:      v.Location,
			Biaya:         v.Fee,
			Kuota:         v.Quota,
			Penyelenggara: "",
			FilledQuota:   int(countParticipant),
		}
	}
	return
}

func (h *TourneyUsecase) CreateTournament(ctx context.Context, req domain.Tournament) (res domain.Tournament, status int, err error) {
	slog.Info("[Usecase][CreateTournament] CreateTournament")
	res, status, err = h.mysqlRepository.CreateTournament(ctx, req)
	if err != nil {
		slog.Error("[Usecase][InquiryTourneyPublic] " + err.Error())
		status = domain.StatusInternalServerError
		return
	}
	return
}
