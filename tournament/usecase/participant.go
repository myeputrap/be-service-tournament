package usecase

import (
	"be-service-tournament/domain"
	"context"
	"log/slog"
	"strconv"
)

func (h *TourneyUsecase) FormPartnershipParticipant(ctx context.Context, req domain.ParticipantDTO) (status int, err error) {
	slog.Info("[Usecase][FormPartnershipParticipant] FormPartnershipParticipant")
	//check if exist
	//TODO Chek if referalCode is right? Also add variable refereal code in request and check it to player 2
	param := make(map[string]string)
	param["id"] = strconv.Itoa(int(req.PlayerTwo))
	user, status, err := h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil {
		slog.Error("[Usecase][CreateAdmin]" + err.Error())
		return
	}
	if *user.ReferalCode != req.ReferalCode {
		status = domain.StatusWrongReferalCode
		err = domain.ErrBadRequest
		return
	}
	param["id"] = strconv.Itoa(int(req.PlayerOne))
	_, status, err = h.mysqlRepository.GetUserByParam(ctx, param)
	if err != nil {
		slog.Error("[Usecase][CreateAdmin]" + err.Error())
		return
	}
	param["id"] = strconv.Itoa(int(req.TournamentID))
	_, status, err = h.mysqlRepository.GetTournamentByParam(ctx, param)
	if err != nil {
		slog.Error("[Usecase][CreateAdmin]" + err.Error())
		return
	}
	isExist, status, err := h.mysqlRepository.IsPlayerExistOnParticipant(ctx, req.TournamentID, req.PlayerOne)
	if err != nil {
		slog.Error("[Usecase][FormPartnershipParticipant] " + err.Error())
		return
	}
	if isExist {
		status = domain.StatusPlayerAlreadyRegistered
		err = domain.ErrBadRequest
		return
	}

	isExist, status, err = h.mysqlRepository.IsPlayerExistOnParticipant(ctx, req.TournamentID, req.PlayerTwo)
	if err != nil {
		slog.Error("[Usecase][FormPartnershipParticipant] " + err.Error())
		return
	}
	if isExist {
		status = domain.StatusPlayerAlreadyRegistered
		err = domain.ErrBadRequest
		return
	}
	//TODO check if tourney is valid?
	//TODO check if partner is the same as gender eligible tournament. For example double male cannot be in mixed tournament
	status, err = h.mysqlRepository.CreateParticipant(ctx, domain.Participant{TournamentID: req.TournamentID, UserAID: req.PlayerOne, UserBID: req.PlayerTwo, State: "Applied"})
	if err != nil {
		slog.Error("[Usecase][FormPartnershipParticipant] " + err.Error())
		return
	}

	return
}

func (h *TourneyUsecase) UpdateParticipant(ctx context.Context, req domain.UpdateParticipantRequest) (status int, err error) {
	slog.Info("[Usecase][UpdateParticipant] UpdateParticipant")
	//checkID
	paramUpdate := make(map[string]string)
	paramGet := make(map[string]string)
	paramGet["id"] = strconv.Itoa(int(req.UserID))
	_, status, err = h.mysqlRepository.GetParticipantByParam(ctx, paramGet)
	if err != nil {
		slog.Error("[Usecase][UpdateParticipant]" + err.Error())
		return
	}
	paramUpdate["state"] = req.Status
	status, err = h.mysqlRepository.UpdateParticipant(ctx, paramUpdate, req.UserID)
	if err != nil {
		slog.Error("[Usecase][UpdateParticipant]" + err.Error())
		return
	}
	return domain.StatusSuccess, nil
}
