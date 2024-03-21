package handler

import (
	"VK_HR/pkg/actorrepo"
	"VK_HR/pkg/filmrepo"
	"VK_HR/pkg/validator"
	validate "github.com/go-playground/validator/v10"
	"net/http"
)

type AppHandler struct {
	ActorsRepo     actorrepo.ActorsRepository
	FilmsRepo      filmrepo.FilmsRepository
	ActorValidator validator.ValueValidator
	FilmValidator  validator.ValueValidator
	validator      *validate.Validate
}

func NewAppHandler(actorsRepo actorrepo.ActorsRepository, filmsRepo filmrepo.FilmsRepository,
	actorValidator, filmValidator validator.ValueValidator) *AppHandler {
	return &AppHandler{
		ActorsRepo:     actorsRepo,
		FilmsRepo:      filmsRepo,
		ActorValidator: actorValidator,
		FilmValidator:  filmValidator,
		validator:      validate.New(),
	}
}

func (handler *AppHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	actor := new(actorrepo.Actor)
	err = read(r, actor)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.validator.Struct(actor)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	insertedID, err := handler.ActorsRepo.Insert(r.Context(), actor)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, insertedResponse{InsertedID: insertedID}); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	updatesParams := r.URL.Query()

	id, err := handler.getIDFromQuery(r.URL.Query(), "actor_id")
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	updatesParams.Del("actor_id")

	setClause, args, err := handler.getUpdateQueryParams(updatesParams, handler.ActorValidator)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.ActorsRepo.Update(r.Context(), setClause, args, id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	id, err := handler.getIDFromQuery(r.URL.Query(), "actor_id")
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.ActorsRepo.Delete(r.Context(), id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	film := new(filmrepo.Film)
	err = read(r, film)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.validator.Struct(film)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	insertedID, err := handler.FilmsRepo.InsertFilm(r.Context(), film)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	err = handler.FilmsRepo.InsertActorsByFilm(r.Context(), insertedID, film.ActorsID)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, insertedResponse{InsertedID: insertedID}); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	updatesParams := r.URL.Query()

	id, err := handler.getIDFromQuery(r.URL.Query(), "film_id")
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	updatesParams.Del("film_id")

	setClause, args, err := handler.getUpdateQueryParams(updatesParams, handler.FilmValidator)
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.FilmsRepo.Update(r.Context(), setClause, args, id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	err := handler.checkAuth(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}

	id, err := handler.getIDFromQuery(r.URL.Query(), "film_id")
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.FilmsRepo.Delete(r.Context(), id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	var name, direction string

	names, exist := values["column_name"]
	if exist {
		name = names[0]
	}

	directions, exist := values["direction"]
	if exist {
		direction = directions[0]
	}

	config, err := filmrepo.NewSortingConfig(validator.ColumnName(name), filmrepo.SortingDirection(direction))
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	films, err := handler.FilmsRepo.SelectBySorting(r.Context(), config)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, films); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) GetFilmByFragment(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	var filmName, actorName string

	filmNames, exist := values["film"]
	if exist {
		filmName = filmNames[0]
	}

	actorNames, exist := values["actor"]
	if exist {
		actorName = actorNames[0]
	}

	films, err := handler.FilmsRepo.SelectByFragment(r.Context(), filmName, actorName)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, films); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	actors, err := handler.ActorsRepo.SelectAll(r.Context())
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err = write(w, actors); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
