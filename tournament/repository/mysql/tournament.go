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

func NewSQLTourneyRepository(dbConn *gorm.DB) domain.SQLTournamentRepository {
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
