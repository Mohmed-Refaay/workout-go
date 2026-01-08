package app

import (
	"backend-go/internals/api"
	"backend-go/internals/store"
	"backend-go/migrations"
	"database/sql"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "Logger: ", log.Ldate|log.Ltime)

	// database
	pgDB, dbErr := store.Open()
	if dbErr != nil {
		return nil, dbErr
	}

	if err := store.MigrateFS(pgDB, migrations.FS, "."); err != nil {
		panic(err)
	}

	// stores
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(userStore, tokenStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HandleHeath(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!\n"))

}
