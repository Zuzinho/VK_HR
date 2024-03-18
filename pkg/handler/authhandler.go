package handler

import (
	"VK_HR/pkg/loginformrepo"
	"VK_HR/pkg/sessionrepo"
	"VK_HR/pkg/userrepo"
	"net/http"
)

type AuthHandler struct {
	LoginFormsRepo loginformrepo.LoginFormsRepository
	UsersRepo      userrepo.UsersRepository
	SessionPacker  sessionrepo.SessionPacker
}

func NewAuthHandler(
	loginFormsRepository loginformrepo.LoginFormsRepository,
	usersRepository userrepo.UsersRepository,
	sessionPacker sessionrepo.SessionPacker,
) *AuthHandler {
	return &AuthHandler{
		LoginFormsRepo: loginFormsRepository,
		UsersRepo:      usersRepository,
		SessionPacker:  sessionPacker,
	}
}

func NewAuthHandlerWithJsonTools(
	loginFormsRepository loginformrepo.LoginFormsRepository,
	usersRepository userrepo.UsersRepository,
	sessionPacker sessionrepo.SessionPacker,
) *AuthHandler {
	return &AuthHandler{
		LoginFormsRepo: loginFormsRepository,
		UsersRepo:      usersRepository,
		SessionPacker:  sessionPacker,
	}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	form := new(loginformrepo.LoginForm)
	err := read(r, form)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	exist, err := handler.LoginFormsRepo.SignIn(r.Context(), form)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if !exist {
		httpError(w, newIncorrectLoginPasswordError(*form), http.StatusBadRequest)
		return
	}

	token, err := handler.SessionPacker.Pack(form.Login)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, authResponse{Token: *token}); err != nil {
		httpError(w, newIncorrectLoginPasswordError(*form), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	form := new(loginformrepo.LoginForm)
	err := read(r, form)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.LoginFormsRepo.SignUp(r.Context(), form)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	err = handler.UsersRepo.Insert(r.Context(), form.Login, userrepo.RegularUser)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	token, err := handler.SessionPacker.Pack(form.Login)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, authResponse{Token: *token}); err != nil {
		httpError(w, newIncorrectLoginPasswordError(*form), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
