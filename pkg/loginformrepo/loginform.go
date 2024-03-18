package loginformrepo

import "context"

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginFormsRepository interface {
	SignUp(ctx context.Context, loginForm *LoginForm) error
	SignIn(ctx context.Context, loginForm *LoginForm) (bool, error)
}
