package main

import (
	"VK_HR/pkg/actorrepo"
	"VK_HR/pkg/env"
	"VK_HR/pkg/filmrepo"
	"VK_HR/pkg/handler"
	"VK_HR/pkg/loginformrepo"
	"VK_HR/pkg/middleware"
	"VK_HR/pkg/router"
	"VK_HR/pkg/sessionrepo"
	"VK_HR/pkg/userrepo"
	"VK_HR/pkg/validator"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Debug("initializing environment vars from .env...")

	env.InitEnv()

	dbConnStr := env.MustDBConnString()
	port := env.MustPort()
	jwtConfig := env.MustJWTConfig()
	maxConnCount := env.MustMaxConnCount()

	log.Debug("opening db connection...")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Debug("checking db connection...")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("up migration...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./tmp/migrate",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Debug("setting max open connections count...")
	db.SetMaxOpenConns(maxConnCount)

	log.Debug("initializing services")
	usersRepo := userrepo.NewUserDBRepository(db)
	actorsRepo := actorrepo.NewActorsDBRepository(db)
	filmsRepo := filmrepo.NewFilmsDBRepository(db)
	sessionManager := sessionrepo.NewSessionManager(jwtConfig)
	loginFormsRepo := loginformrepo.NewLoginFormsDBRepository(db)

	actorsValidator := validator.NewActorsValidator()
	filmsValidator := validator.NewFilmsValidator()

	authHandler := handler.NewAuthHandler(loginFormsRepo, usersRepo, sessionManager)
	appHandler := handler.NewAppHandler(actorsRepo, filmsRepo, actorsValidator, filmsValidator)

	middle := middleware.NewMiddleware(sessionManager, usersRepo)
	rout := router.NewRouter(middle)

	rout.HandleFunc(router.NewPair("/login", "POST"), authHandler.Login)
	rout.HandleFunc(router.NewPair("/register", "POST"), authHandler.Register)

	rout.HandleFuncWithAuth(router.NewPair("/actor", "POST"), appHandler.AddActor)
	rout.HandleFuncWithAuth(router.NewPair("/actor", "PUT"), appHandler.UpdateActor)
	rout.HandleFuncWithAuth(router.NewPair("/actor", "DELETE"), appHandler.DeleteActor)
	rout.HandleFunc(router.NewPair("/actors", "GET"), appHandler.GetAllActors)

	rout.HandleFuncWithAuth(router.NewPair("/film", "POST"), appHandler.AddFilm)
	rout.HandleFuncWithAuth(router.NewPair("/film", "PUT"), appHandler.UpdateFilm)
	rout.HandleFuncWithAuth(router.NewPair("/film", "DELETE"), appHandler.DeleteFilm)
	rout.HandleFunc(router.NewPair("/films", "GET"), appHandler.GetAllFilms)
	rout.HandleFunc(router.NewPair("/film", "GET"), appHandler.GetFilmByFragment)

	api := rout.PackInMiddleware()

	log.Debugf("started on port :%s", port)

	log.Debug(http.ListenAndServe(":"+port, api))
}
