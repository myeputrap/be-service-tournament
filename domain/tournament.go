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

type TournamentUsecase interface {
	Login(ctx context.Context, req RequestLogin) (r map[string]interface{}, status int, err error)
}

type SQLTournamentRepository interface {
	GetUserLogin(ctx context.Context, req RequestLogin) (user User, status int, err error)
}
