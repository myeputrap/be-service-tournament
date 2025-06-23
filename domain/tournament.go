package domain

import "context"

type RequestLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
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

type TournamentUsecase interface {
	Login(ctx context.Context, req RequestLogin) (r map[string]interface{}, status int, err error)
	InquiryTourneyPublic(ctx context.Context) (response []InquiryTourneyPublicResponse, status int, err error)
}

type SQLTournamentRepository interface {
	GetUserLogin(ctx context.Context, req RequestLogin) (user User, status int, err error)
	GetTournament(ctx context.Context) (response []Tournament, status int, err error)
	CountParticipantTournament(ctx context.Context, tourneyID int32) (count int64, status int, err error)
}
