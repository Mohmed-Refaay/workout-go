package routes

import (
	"backend-go/internals/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.Authenticate)
		r.Use(app.UserMiddleware.RequireUser)

		// Workout
		r.Get("/workouts/{id}", app.WorkoutHandler.GetWorkoutById)
		r.Delete("/workouts/{id}", app.WorkoutHandler.DeleteWorkoutById)
		r.Put("/workouts/{id}", app.WorkoutHandler.UpdateWorkout)
		r.Post("/workouts", app.WorkoutHandler.CreateWorkout)
	})

	// User
	r.Post("/users", app.UserHandler.CreateUser)
	r.Post("/token/authenticate", app.TokenHandler.Authenticate)

	return r
}
