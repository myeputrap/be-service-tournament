package usecase

import (
	"be-service-tournament/domain"
	"context"
	"crypto/md5"
	"encoding/hex"
	"log/slog"
	"strconv"
	"time"

	"github.com/spf13/viper"
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

func (h *TourneyUsecase) UpdateParticipant(ctx context.Context, req domain.UpdateParticipantStatusRequest) (status int, err error) {
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

func (h *TourneyUsecase) CreatePaymentProofImage(ctx context.Context, req domain.RequestPaymentProffImage) (status int, err error) {
	slog.Info("[Usecase][CreatePaymentProofImage] CreatePaymentProofImage")
	assetpath := viper.GetString("server.http.asset_path")
	param := make(map[string]string)
	param["id"] = strconv.Itoa(int(req.ParticipantID))
	//check participantID
	participant, status, err := h.mysqlRepository.GetParticipantByParam(ctx, param) //TODO ask which statusregistered allow to upload image paymentProof
	if err != nil {
		slog.Error("[Usecase][CreatePaymentProofImage]" + err.Error())
		return
	}
	fn := strconv.Itoa(int(participant.ID)) + strconv.FormatInt(time.Now().Unix(), 10)
	hash := md5.Sum([]byte(fn))
	encoded := hex.EncodeToString(hash[:])

	fileName, err := h.assetRepository.SaveFile(req.Images, assetpath+"/images/payment_proof", encoded)
	if err != nil {
		slog.Error("[Usecase][CreateTheme][SaveFileBackground]", "Err", err)
		status = domain.StatusInternalServerError
		return
	}
	status, err = h.mysqlRepository.DynamicEditTable(ctx, map[string]string{"payment_proof": fileName}, int(participant.ID), &domain.Participant{})
	if err != nil {
		slog.Error("[Usecase][DynamicEditTable][DynamicEditTable]", "Err", err)
		status = domain.StatusInternalServerError
		return
	}
	return
}

func (h *TourneyUsecase) GetAllParticipant(ctx context.Context, req domain.GetAllParticipantRequest) (res domain.GetAllParticipantResponse, status int, err error) {
	slog.Info("[Usecase][GetParticipantList] GetParticipantList")
	req.Offset = (req.Page - 1) * req.Limit
	out, count, status, err := h.mysqlRepository.GetAllParticipant(ctx, req)
	if err != nil {
		slog.Error("[Usecase][GetParticipantList] " + err.Error())
		return
	}

	res.Metadata = domain.MetaData{
		TotalData: uint(count),
		TotalPage: (uint(count) + uint(req.Limit) - 1) / uint(req.Limit),
		Page:      uint(req.Page),
		Limit:     uint(req.Limit),
		Sort:      req.Sort,
		Order:     req.Order,
	}
	res.Data = out

	return res, domain.StatusSuccess, nil
}
