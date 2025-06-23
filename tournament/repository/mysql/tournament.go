package mysql

import (
	"be-service-tournament/domain"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type tourneyMySQLRepository struct {
	Conn *gorm.DB
}

func NewSQLTournamentRepository(dbConn *gorm.DB) domain.SQLTournamentRepository {
	return &tourneyMySQLRepository{
		Conn: dbConn,
	}
}

func (t *tourneyMySQLRepository) GetUserLogin(ctx context.Context, req domain.RequestLogin) (user domain.User, status int, err error) {
	slog.Info("[Repository][Login] Login")

	var response domain.User
	query := t.Conn.WithContext(ctx)

	query = query.Where("email = ?", req.Email)

	err = query.First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = domain.StatusNotFound
			err = domain.ErrNotFound
			return
		}
		slog.Error("[Repository][GetUserLogin] Row scan error", "Err", err.Error())
		status = domain.StatusInternalServerError
		return
	}
	return response, domain.StatusSuccess, nil
}
func (t *tourneyMySQLRepository) GetTournament(ctx context.Context) (responses []domain.Tournament, status int, err error) {
	slog.Info("[Repository][GetTournament] GetTournament")
	db := t.Conn.WithContext(ctx).Model(&domain.Tournament{}).Where("NOW() BETWEEN start_date and end_date")
	result := db.Find(&responses)
	if result.Error != nil {
		slog.Error("[Repository][GetTournament] err", "", result.Error)
		status = domain.StatusInternalServerError
		return

	}
	if result.RowsAffected == 0 {
		status = domain.StatusNotFound
		err = domain.ErrNotFound
		slog.Error("[Repository][GetTournament] GetTournament not found", "Err", err.Error())
		return
	}
	return
}

func (t *tourneyMySQLRepository) CountParticipantTournament(ctx context.Context, tourneyID int32) (count int64, status int, err error) {
	slog.Info("[Repository][CountParticipantTournament] CountParticipantTournament")

	err = t.Conn.WithContext(ctx).
		Model(&domain.Participant{}).
		Count(&count).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return
	}

	return
}
