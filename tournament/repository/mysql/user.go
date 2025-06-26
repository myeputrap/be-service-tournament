package mysql

import (
	"be-service-tournament/domain"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

func (t *tourneyMySQLRepository) GetUserByParam(ctx context.Context, param map[string]string) (res *domain.User, status int, err error) {
	slog.Info("[Repository][GetUserByParam] GetUserByParam")
	var ttB domain.User
	query := t.Conn.WithContext(ctx)

	for column, value := range param {
		query = query.Where(column+" = ?", value)
	}

	err = query.First(&ttB).Error
	if err != nil {
		slog.Error("[Repository][GetUserByParam] err", "", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.StatusNotFound, domain.ErrNotFound
		}
		status = domain.StatusInternalServerError
		return
	}
	return &ttB, domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) CreateUser(ctx context.Context, req domain.User) (res domain.User, status int, err error) {
	slog.Info("[Repository][CreateUser] CreateUser")
	result := t.Conn.WithContext(ctx).Create(&req)
	err = result.Error
	if err != nil {
		slog.Error("[Repository][CreateUser] err", "", err)
		status = domain.StatusInternalServerError
		return
	}
	res = req
	status = domain.StatusSuccessCreate
	return
}

func (t *tourneyMySQLRepository) GetUserPartner(ctx context.Context, req domain.GetAllUserRequestPartner) (res []domain.User, count int, status int, err error) {
	slog.Info("[Repository][GetUserPartner] GetUserPartner")
	db := t.Conn.WithContext(ctx).Model(&domain.User{}).Where("role_id != 2").Where("id != ?", req.UIDSearcher)
	if req.Gender != nil && *req.Gender != "" {
		db = db.Where("gender = ?", req.Gender)
	}
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE %?%", req.Name)
	}
	if req.Tier != nil && *req.Tier != "" {
		db = db.Where("tier = ?", req.Tier)
	}
	result := db.Find(&res)
	if result.Error != nil {
		slog.Error("[Repository][GetUserPartner] err", "", result.Error)
		status = domain.StatusInternalServerError
		return
	}

	if result.RowsAffected == 0 {
		status = domain.StatusNotFound
		err = domain.ErrNotFound
		slog.Error("[Repository][GetUserPartner] GetUserPartner not found", "Err", err.Error())
		return
	}
	return
}
