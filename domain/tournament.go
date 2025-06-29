package domain

import (
	"context"
	"mime/multipart"

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
	ID           int64   `json:"id"`
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
	Count int              `json:"count"`
	Data  []UserPartnerDTO `json:"data"`
}
type ParticipantDTO struct {
	TournamentID int64  `json:"tournament_id" validate:"required"`
	PlayerOne    int64  `json:"player_one" validate:"required"`
	PlayerTwo    int64  `json:"player_two" validate:"required"`
	ReferalCode  string `json:"referal_code" validate:"required"`
}

type UpdateParticipantStatusRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type UserDTO struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Gender      string `json:"gender"`
	TierName    string `json:"tier_name"`
}

type GetAllTournamentRequest struct {
	Page   int64  `json:"page" validate:"required"`
	Limit  int64  `json:"limit" validate:"required"`
	Offset int64  `json:"offeset"`
	Sort   string `json:"sort" validate:"required"`
	Order  string `json:"order" validate:"required"`
}
type GetTournamentParticipantRequest struct {
	TournamentID int64 `json:"tournament_id" `
}
type GetAllTournamentResponse struct {
	Metadata MetaData     `json:"meta"`
	Data     []Tournament `json:"data"`
}
type RequestPaymentProffImage struct {
	ParticipantID int64                 `json:"participant_id"`
	Images        *multipart.FileHeader `json:"images" form:"images" validate:"required"`
}

type GetAllParticipantRequest struct {
	Page   int64  `json:"page" validate:"required"`
	Limit  int64  `json:"limit" validate:"required"`
	Offset int64  `json:"offeset"`
	Sort   string `json:"sort" validate:"required"`
	Order  string `json:"order" validate:"required"`
}

type GetAllParticipantResponse struct {
	Metadata MetaData      `json:"meta"`
	Data     []Participant `json:"data"`
}
type TournamentUsecase interface {
	Login(ctx context.Context, req RequestLogin) (r map[string]interface{}, status int, err error)
	InquiryTourneyPublic(ctx context.Context) (response []InquiryTourneyPublicResponse, status int, err error)

	CreateUser(ctx context.Context, req UserRequestDTO) (res *User, multiErr *MultipleErrorResponse, status int, err error)
	GetUserByID(ctx context.Context, id int64) (res *User, status int, err error)
	GetUserPartner(ctx context.Context, req GetAllUserRequestPartner) (res GetUserPartnerResponseDTO, status int, err error)
	CreateTournament(ctx context.Context, req Tournament) (res Tournament, status int, err error)
	FormPartnershipParticipant(ctx context.Context, req ParticipantDTO) (status int, err error)
	CreateAdmin(ctx context.Context, req UserRequestDTO) (res *User, multiErr *MultipleErrorResponse, status int, err error)
	UpdateParticipant(ctx context.Context, req UpdateParticipantStatusRequest) (status int, err error)
	GetAllTournament(ctx context.Context, req GetAllTournamentRequest) (res GetAllTournamentResponse, status int, err error)
	GetAllParticipant(ctx context.Context, req GetAllParticipantRequest) (res GetAllParticipantResponse, status int, err error)
	CreatePaymentProofImage(ctx context.Context, req RequestPaymentProffImage) (status int, err error)
	GetTournamentParticipant(ctx context.Context, req GetTournamentParticipantRequest) (res []UserDTO, status int, err error)
	GetUserTournament(ctx context.Context, id int64) (res []Tournament, status int, err error)
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
	GetTournamentByIDs(ctx context.Context, param []int64) (res []Tournament, status int, err error)
	GetParticipantByParam(ctx context.Context, param map[string]string) (res *Participant, status int, err error)
	GetParticipantByParamArray(ctx context.Context, param map[string]string) (res []Participant, status int, err error)
	UpdateParticipant(ctx context.Context, params map[string]string, id int64) (status int, err error)
	IsPlayerExistOnParticipant(ctx context.Context, tourneyID int64, userID int64) (isExist bool, status int, err error)
	CreateParticipant(ctx context.Context, req Participant) (status int, err error)
	GetAllTournament(ctx context.Context, req GetAllTournamentRequest) (res []Tournament, count int64, status int, err error)
	GetAllParticipant(ctx context.Context, req GetAllParticipantRequest) (res []Participant, count int64, status int, err error)
	DynamicEditTable(ctx context.Context, params map[string]string, id int, model any) (status int, err error)
	GetTournamentParticipant(ctx context.Context, req []int64) (res []User, status int, err error)
}

type AssetRepository interface {
	Remove(path string, generatedFileName string) (err error)
	SaveFile(asset *multipart.FileHeader, path string, generatedFileName string) (fileName string, err error)
}
