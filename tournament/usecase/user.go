package usecase

import (
	"be-service-tournament/domain"
	"be-service-tournament/helper"
	"context"
	"log/slog"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (h *TourneyUsecase) CreateUser(ctx context.Context, req domain.UserRequestDTO) (respon *domain.User, multiError *domain.MultipleErrorResponse, status int, err error) {
	slog.Info("[Usecase][CreateUser] CreateUser")
	param := make(map[string]string)
	param["email"] = req.Email

	user, status, err := h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil && status != domain.StatusNotFound {
		slog.Error("[Usecase][CreateUser]" + err.Error())
		return
	}
	multiErr := &domain.MultipleErrorResponse{}
	if user != nil {
		multiErr.Email = append(multiErr.Email, domain.ErrorResponseMsg{Error: "email already exists"})
	}

	param = make(map[string]string)
	param["username"] = req.Username

	user, status, err = h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil && status != domain.StatusNotFound {
		slog.Error("[Usecase][CreateUser]" + err.Error())
		return
	}
	if user != nil {
		multiErr.Username = append(multiErr.Username, domain.ErrorResponseMsg{Error: "username already exist"})
	}

	if req.Password != req.ConfirmPassword {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password and confirm password not match"})
	}
	if len(req.Password) < 8 {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password length must be more than 8"})
	}
	if !helper.IsContainNumber(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain number"})
	}
	if !helper.IsContainCapital(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain capital letter"})
	}
	if !helper.IsContainSpecialCharacter(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain special character"})
	}

	if len(multiErr.Email) > 0 || len(multiErr.Username) > 0 || len(multiErr.Password) > 0 {
		return nil, multiErr, domain.StatusBadRequest, domain.ErrBadRequest
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		status = domain.StatusInternalServerError
		slog.Error("[Usecase][CreateUser] Generate password", "Err", err.Error())
		return
	}
	var tierID int64 = 1
	res, status, err := h.mysqlRepository.CreateUser(ctx, domain.User{
		Email:          req.Email,
		Username:       req.Username,
		HashedPassword: string(password),
		PhoneNumber:    req.PhoneNumber,
		FullName:       req.FullName,
		Gender:         req.Gender,
		RoleID:         2,
		TierID:         &tierID,
	})
	if err != nil {
		slog.Error("[Usecase][Login]" + err.Error())
		return
	}

	return &res, nil, domain.StatusSuccessCreate, nil
}

func (h *TourneyUsecase) CreateAdmin(ctx context.Context, req domain.UserRequestDTO) (respon *domain.User, multiError *domain.MultipleErrorResponse, status int, err error) {
	slog.Info("[Usecase][CreateAdmin] CreateAdmin")
	param := make(map[string]string)
	param["email"] = req.Email

	user, status, err := h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil && status != domain.StatusNotFound {
		slog.Error("[Usecase][CreateAdmin]" + err.Error())
		return
	}
	multiErr := &domain.MultipleErrorResponse{}
	if user != nil {
		multiErr.Email = append(multiErr.Email, domain.ErrorResponseMsg{Error: "email already exists"})
	}

	param = make(map[string]string)
	param["username"] = req.Username

	user, status, err = h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil && status != domain.StatusNotFound {
		slog.Error("[Usecase][CreateAdmin]" + err.Error())
		return
	}
	if user != nil {
		multiErr.Username = append(multiErr.Username, domain.ErrorResponseMsg{Error: "username already exist"})
	}

	if req.Password != req.ConfirmPassword {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password and confirm password not match"})
	}
	if len(req.Password) < 8 {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password length must be more than 8"})
	}
	if !helper.IsContainNumber(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain number"})
	}
	if !helper.IsContainCapital(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain capital letter"})
	}
	if !helper.IsContainSpecialCharacter(req.Password) {
		multiErr.Password = append(multiErr.Password, domain.ErrorResponseMsg{Error: "password must be contain special character"})
	}

	if len(multiErr.Email) > 0 || len(multiErr.Username) > 0 || len(multiErr.Password) > 0 {
		return nil, multiErr, domain.StatusBadRequest, domain.ErrBadRequest
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		status = domain.StatusInternalServerError
		slog.Error("[Usecase][CreateAdmin] Generate password", "Err", err.Error())
		return
	}
	referalCode := helper.UuidGenerator()
	res, status, err := h.mysqlRepository.CreateUser(ctx, domain.User{
		Email:          req.Email,
		Username:       req.Username,
		HashedPassword: string(password),
		PhoneNumber:    req.PhoneNumber,
		FullName:       req.FullName,
		Gender:         req.Gender,
		RoleID:         1,
		ReferalCode:    &referalCode,
	})
	if err != nil {
		slog.Error("[Usecase][Login]" + err.Error())
		return
	}

	return &res, nil, domain.StatusSuccessCreate, nil
}

func (h *TourneyUsecase) GetUserPartner(ctx context.Context, req domain.GetAllUserRequestPartner) (res domain.GetUserPartnerResponseDTO, status int, err error) {
	slog.Info("[Usecase][GetUserPartner] GetUserPartner")
	user, _, status, err := h.mysqlRepository.GetUserPartner(ctx, req)
	if err != nil {
		slog.Error("[Usecase][GetUserPartner]" + err.Error())
		return
	}
	res.Count = len(user)
	res.Data = make([]domain.UserPartnerDTO, len(user))
	for i, v := range user {
		res.Data[i] = domain.UserPartnerDTO{
			ID:           v.ID,
			Email:        v.Email,
			Username:     v.Username,
			PhoneNumber:  v.PhoneNumber,
			FullName:     v.FullName,
			Gender:       v.Gender,
			Tier:         "",
			PhotoProfile: v.PhotoProfile,
		}
	}
	return res, domain.StatusSuccess, nil
}

func (h *TourneyUsecase) GetUserByID(ctx context.Context, id int64) (res *domain.User, status int, err error) {
	slog.Info("[Usecase][GetUserByID] GetUserByID")
	param := make(map[string]string)
	param["id"] = strconv.Itoa(int(id))
	res, status, err = h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil {
		slog.Error("[Usecase][GetUserByID]" + err.Error())
		return
	}
	return
}

func (h *TourneyUsecase) GetUserTournament(ctx context.Context, id int64) (res []domain.Tournament, status int, err error) {
	slog.Info("[Usecase][GetUserTournament] GetUserTournament")
	//check userID
	param := make(map[string]string)
	param["id"] = strconv.Itoa(int(id))
	_, status, err = h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil {
		slog.Error("[Usecase][GetUserTournament] " + err.Error())
		return
	}

	//getTournamentID
	var idArray []int64
	param = make(map[string]string)
	param["user_a_id"] = strconv.Itoa(int(id))
	param["user_b_id"] = strconv.Itoa(int(id))
	participant, status, err := h.mysqlRepository.GetParticipantByParamArray(ctx, param)
	if err != nil {
		slog.Error("[Usecase][GetUserTournament] " + err.Error())
		return
	}
	for _, v := range participant {
		idArray = append(idArray, v.TournamentID)
	}
	tourney, status, err := h.mysqlRepository.GetTournamentByIDs(ctx, idArray)
	if err != nil {
		slog.Error("[Usecase][GetUserTournament] " + err.Error())
		return
	}
	return tourney, domain.StatusSuccess, nil
}
