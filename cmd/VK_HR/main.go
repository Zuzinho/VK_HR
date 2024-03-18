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
	"log"
	"net/http"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func main() {
	env.InitEnv()

	dbConnStr := env.MustDBConnString()
	port := env.MustPort()
	jwtConfig := env.MustJWTConfig()

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../tmp/migrate",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)

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
	rout.HandleFuncWithAuth(router.NewPair("/actor", "UPDATE"), appHandler.UpdateActor)
	rout.HandleFuncWithAuth(router.NewPair("/actor", "DELETE"), appHandler.DeleteActor)
	rout.HandleFuncWithAuth(router.NewPair("/actors", "GET"), appHandler.GetAllActors)

	rout.HandleFuncWithAuth(router.NewPair("/film", "POST"), appHandler.AddFilm)
	rout.HandleFuncWithAuth(router.NewPair("/film", "UPDATE"), appHandler.UpdateFilm)
	rout.HandleFuncWithAuth(router.NewPair("/film", "DELETE"), appHandler.DeleteFilm)
	rout.HandleFuncWithAuth(router.NewPair("/films", "GET"), appHandler.GetAllFilms)
	rout.HandleFuncWithAuth(router.NewPair("/film", "GET"), appHandler.GetFilmByFragment)

	api := rout.PackInMiddleware()

	log.Printf("started on port :%s", port)

	log.Println(http.ListenAndServe(":"+port, api))
}
