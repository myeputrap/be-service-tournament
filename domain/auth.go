package domain

type AuthRequest struct {
	Method string
	Path   string
	Token  string
}
