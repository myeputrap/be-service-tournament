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
