package mysql

import (
	"be-service-tournament/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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
func (t *tourneyMySQLRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := t.Conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
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

func (t *tourneyMySQLRepository) DeleteTableBasedOnParameter(ctx context.Context, table string, param string, id int64, tx *gorm.DB) (status int, err error) {
	slog.Info("[Repository][DeleteTableBasedOnParameter] DeleteTableBasedOnParameter")
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, param)
	slog.Debug("[Repository][DeleteTableBasedOnParameter] Query", ":", query+"_"+table+"_"+param+"_"+strconv.FormatInt(id, 10))

	result := tx.WithContext(ctx).Exec(query, id)
	if result.Error != nil {
		slog.Error("[Repository][DeleteTableBasedOnParameter] Exec error", "", result.Error)
		return domain.StatusInternalServerError, result.Error
	}

	affected := result.RowsAffected
	if affected == 0 {
		err = errors.New("not found")
		slog.Error("[Repository][DeleteTableBasedOnParameter] err", "", err)
		return domain.StatusNotFound, err
	}

	return domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) CreateTournament(ctx context.Context, req domain.Tournament) (res domain.Tournament, status int, err error) {
	slog.Info("[Repository][CreateTournament] CreateTournament")
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

func (t *tourneyMySQLRepository) GetTournamentByParam(ctx context.Context, param map[string]string) (res *domain.Tournament, status int, err error) {
	slog.Info("[Repository][GetTournamentByParam] GetTournamentByParam")
	var ttB domain.Tournament
	query := t.Conn.WithContext(ctx)

	for column, value := range param {
		query = query.Where(column+" = ?", value)
	}

	err = query.First(&ttB).Error
	if err != nil {
		slog.Error("[Repository][GetTournamentByParam] err", "", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.StatusNotFound, domain.ErrNotFound
		}
		status = domain.StatusInternalServerError
		return
	}
	return &ttB, domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) GetAllTournament(ctx context.Context, req domain.GetAllTournamentRequest) (res []domain.Tournament, count int64, status int, err error) {
	db := t.Conn.WithContext(ctx).Model(&domain.Tournament{}).Order(fmt.Sprintf("%s %s", req.Sort, req.Order)).Limit(int(req.Limit)).Offset(int(req.Offset))
	err = db.Count(&count).Error
	if err != nil {
		slog.Error("[Repository][GetAllTournament] err", "", err.Error())
		status = domain.StatusInternalServerError
		return
	}
	result := db.Find(&res)
	if result.Error != nil {
		slog.Error("[Repository][GetAllTournament] err", "", result.Error)
		status = domain.StatusInternalServerError
		return
	}

	if result.RowsAffected == 0 {
		status = domain.StatusNotFound
		err = domain.ErrNotFound
		slog.Error("[Repository][GetAllTournament] GetAllTournament not found", "Err", err.Error())
		return
	}
	return
}

func (r *tourneyMySQLRepository) DynamicEditTable(ctx context.Context, params map[string]string, id int, model any) (status int, err error) {
	slog.Info("[Repository][DynamicEditTable] Update Table")

	updateData := map[string]interface{}{}
	for key, value := range params {
		if strings.ToLower(value) == "null" {
			updateData[key] = nil
		} else {
			updateData[key] = value
		}
	}

	err = r.Conn.WithContext(ctx).
		Model(model).
		Where("id = ?", id).
		Updates(updateData).Error

	if err != nil {
		slog.Error("[Repository][DynamicEditTable] Update error", slog.Any("error", err))
		return domain.StatusInternalServerError, err
	}

	return domain.StatusSuccess, nil
}

func (t *tourneyMySQLRepository) GetTournamentByIDs(ctx context.Context, param []int64) (res []domain.Tournament, status int, err error) {
	db := t.Conn.WithContext(ctx).Model(&domain.Tournament{}).Where("id in ?", param)

	result := db.Find(&res)
	if result.Error != nil {
		slog.Error("[Repository][GetAllTournament] err", "", result.Error)
		status = domain.StatusInternalServerError
		return
	}

	if result.RowsAffected == 0 {
		status = domain.StatusNotFound
		err = domain.ErrNotFound
		slog.Error("[Repository][GetAllTournament] GetAllTournament not found", "Err", err.Error())
		return
	}
	return res, domain.StatusSuccess, nil
}
