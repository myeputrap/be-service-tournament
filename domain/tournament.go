package domain

import (
	"context"

	"gorm.io/gorm"
)

type RequestLogin struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
type MetaData struct {
	TotalData uint   `json:"total_data"`
	TotalPage uint   `json:"total_page"`
	Page      uint   `json:"page"`
	Limit     uint   `json:"limit"`
	Sort      string `json:"sort"`
	Order     string `json:"order"`
}
type InquiryTourneyPublicResponse struct {
	Name          string  `json:"name"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Location      string  `json:"location"`
	Format        string  `json:"format"`
	Biaya         float64 `json:"biaya"`
	Kuota         int     `json:"kuota"`
	Penyelenggara string  `json:"penyelanggara"`
	FilledQuota   int     `json:"kuota_terisi"`
}
type MultipleErrorResponse struct {
	Email    []ErrorResponseMsg `json:"email,omitempty"`
	Username []ErrorResponseMsg `json:"username,omitempty"`
	Password []ErrorResponseMsg `json:"password,omitempty"`
}

type ErrorResponseMsg struct {
	Error string `json:"error,omitempty"`
}

type UserPartnerDTO struct {
	Email        string  `json:"email"`
	Username     string  `json:"username"`
	PhoneNumber  string  `json:"phone_number"`
	FullName     string  `json:"full_name"`
	Gender       string  `json:"gender"`
	Tier         string  `json:"tier"`
	PhotoProfile *string `json:"photo_profile"`
}

type GetAllUserRequestPartner struct {
	Gender      *string `json:"gender"`
	Name        *string `json:"name"`
	Tier        *string `json:"tier"`
	UIDSearcher int64
}
type GetUserPartnerResponseDTO struct {
	Count int              `json:"meta"`
	Data  []UserPartnerDTO `json:"data"`
}

type TournamentUsecase interface {
	Login(ctx context.Context, req RequestLogin) (r map[string]interface{}, status int, err error)
	InquiryTourneyPublic(ctx context.Context) (response []InquiryTourneyPublicResponse, status int, err error)
	CreateUser(ctx context.Context, req UserRequestDTO) (res *User, multiErr *MultipleErrorResponse, status int, err error)
	GetUserPartner(ctx context.Context, req GetAllUserRequestPartner) (res GetUserPartnerResponseDTO, status int, err error)
	CrateTournament(ctx context.Context, req Tournament) (res Tournament, status int, err error)
}

type SQLTournamentRepository interface {
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	GetUserLogin(ctx context.Context, req RequestLogin) (user User, status int, err error)
	GetTournament(ctx context.Context) (response []Tournament, status int, err error)
	CountParticipantTournament(ctx context.Context, tourneyID int32) (count int64, status int, err error)
	DeleteTableBasedOnParameter(ctx context.Context, table string, param string, id int64, tx *gorm.DB) (status int, err error)
	CreateUser(ctx context.Context, req User) (res User, status int, err error)
	GetUserByParam(ctx context.Context, param map[string]string) (res *User, status int, err error)
	GetUserPartner(ctx context.Context, req GetAllUserRequestPartner) (res []User, count int, status int, err error)
	CreateTournament(ctx context.Context, req Tournament) (res Tournament, status int, err error)
	GetTournamentByParam(ctx context.Context, param map[string]string) (res *Tournament, status int, err error)
	GetParticipantByParam(ctx context.Context, param map[string]string) (res *Participant, status int, err error)
	UpdateParticipant(ctx context.Context, params map[string]string, id int64) (status int, err error)
}
