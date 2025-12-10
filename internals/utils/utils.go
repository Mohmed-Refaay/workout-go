package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func GetParamId(r *http.Request) (int64, error) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		return 0, fmt.Errorf("id is not provided")
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 62)
	if err != nil {
		return 0, err
	}

	return (workoutId), nil
}
