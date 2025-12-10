package api

import (
	"encoding/json"
	"log"
	"net/http"

	"backend-go/internals/store"
	"backend-go/internals/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: store,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.GetParamId(r)
	if err != nil {
		wh.logger.Printf("Error: GetParamId %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "id is missing"})
		return
	}

	wo, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("Error: GetWorkoutById %v\n", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"data": wo,
	}); err != nil {
		wh.logger.Printf("Error: WriteJson %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to encode response"})
		return
	}
}

func (wh *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.GetParamId(r)
	if err != nil {
		wh.logger.Printf("Error: GetParamId %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "id is missing"})
		return
	}

	existingWo, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("Error: GetWorkoutById %v\n", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
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
		wh.logger.Printf("Error: Decode request body %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
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
		wh.logger.Printf("Error: UpdateWorkout %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update workout"})
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, utils.Envelope{"data": newWo}); err != nil {
		wh.logger.Printf("Error: WriteJson %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to encode response"})
		return
	}
}
func (wh *WorkoutHandler) DeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.GetParamId(r)
	if err != nil {
		wh.logger.Printf("Error: GetParamId %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "id is missing"})
		return
	}

	err = wh.workoutStore.DeleteWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("Error: DeleteWorkoutById %v\n", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, utils.Envelope{"message": "workout deleted successfully"}); err != nil {
		wh.logger.Printf("Error: WriteJson %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to encode response"})
		return
	}
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var wo store.Workout
	if err := json.NewDecoder(r.Body).Decode(&wo); err != nil {
		wh.logger.Printf("Error: Decode request body %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	newWo, err := wh.workoutStore.CreateWorkout(&wo)
	if err != nil {
		wh.logger.Printf("Error: CreateWorkout %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	if err := utils.WriteJson(w, http.StatusCreated, utils.Envelope{"data": newWo}); err != nil {
		wh.logger.Printf("Error: WriteJson %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to encode response"})
		return
	}
}
