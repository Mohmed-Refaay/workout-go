package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"backend-go/internals/store"
	"backend-go/internals/utils"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: store,
	}
}

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 62)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	wo, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, wo); err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong!", http.StatusInternalServerError)
		return
	}
}

func (wh *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 62)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	existingWo, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong!", http.StatusInternalServerError)
		return
	}

	type UpdateWorkoutScheme struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	updatedWo := &UpdateWorkoutScheme{}
	if err := json.NewDecoder(r.Body).Decode(updatedWo); err != nil {
		http.Error(w, "Bad input", http.StatusBadRequest)
		return
	}

	if updatedWo.Title != nil {
		existingWo.Title = *updatedWo.Title
	}
	if updatedWo.Description != nil {
		existingWo.Description = *updatedWo.Description
	}
	if updatedWo.DurationMinutes != nil {
		existingWo.DurationMinutes = *updatedWo.DurationMinutes
	}
	if updatedWo.CaloriesBurned != nil {
		existingWo.CaloriesBurned = *updatedWo.CaloriesBurned
	}
	if updatedWo.Entries != nil {
		existingWo.Entries = updatedWo.Entries
	}

	newWo, err := wh.workoutStore.UpdateWorkout(existingWo)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong!", http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJson(w, http.StatusCreated, newWo); err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong!", http.StatusInternalServerError)
		return
	}
}
func (wh *WorkoutHandler) DeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 62)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	err = wh.workoutStore.DeleteWorkoutById(workoutId)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Deleted correctly")
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var wo store.Workout
	if err := json.NewDecoder(r.Body).Decode(&wo); err != nil {
		fmt.Println(err)
		http.Error(w, "Wrong data type", http.StatusInternalServerError)
		return
	}

	newWo, err := wh.workoutStore.CreateWorkout(&wo)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add data", http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJson(w, http.StatusCreated, newWo); err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong!", http.StatusInternalServerError)
		return
	}
}
