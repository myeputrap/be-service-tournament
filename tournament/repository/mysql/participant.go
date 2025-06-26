package mysql

import (
	"be-service-tournament/domain"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

func (t *tourneyMySQLRepository) GetParticipantByParam(ctx context.Context, param map[string]string) (res *domain.Participant, status int, err error) {
	slog.Info("[Repository][GetParticipantByParam] GetParticipantByParam")
	var ttB domain.Participant
	query := t.Conn.WithContext(ctx)

	for column, value := range param {
		query = query.Where(column+" = ?", value)
	}

	err = query.First(&ttB).Error
	if err != nil {
		slog.Error("[Repository][GetParticipantByParam] err", "", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.StatusNotFound, domain.ErrNotFound
		}
		status = domain.StatusInternalServerError
		return
	}
	return &ttB, domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) UpdateParticipant(ctx context.Context, params map[string]string, id int64) (status int, err error) {
	slog.Info("[Repository][UpdateParticipant] UpdateParticipant")
	tx := t.Conn.WithContext(ctx).Model(&domain.Participant{})

	values := make(map[string]interface{})

	for key, value := range params {
		values[key] = value
	}

	err = tx.Where("id = ?", id).Updates(values).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = domain.StatusNotFound
			return
		}
		status = domain.StatusInternalServerError
		return
	}
	status = domain.StatusSuccess
	return
}

func (t *tourneyMySQLRepository) IsPlayerExistOnParticipant(ctx context.Context, tourneyID int64, userID int64) (isExist bool, status int, err error) {
	slog.Info("[Repository][IsPlayerExistOnParticipant] IsPlayerExistOnParticipant")
	var count int64

	err = t.Conn.Model(&domain.Participant{}).
		Where("tournament_id = ? AND (user_a_id = ? OR user_b_id = ?)", tourneyID, userID, userID).
		Count(&count).Error

	if err != nil {
		status = domain.StatusInternalServerError
		return
	}

	if count > 0 {
		isExist = true
		return
	}
	return false, domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) CreateParticipant(ctx context.Context, req domain.Participant) (status int, err error) {
	slog.Info("[Repository][CreateParticipant] CreateParticipant")
	result := t.Conn.WithContext(ctx).Create(&req)
	err = result.Error
	if err != nil {
		slog.Error("[Repository][CreateParticipant] err", "", err)
		status = domain.StatusInternalServerError
		return
	}
	status = domain.StatusSuccessCreate
	return
}
