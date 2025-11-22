package store

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	ID              int64          `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID        int64  `json:"id"`
	WorkoutId int64  `json:"workout_id"`
	Name      string `json:"name"`
	Notes     string `json:"notes"`

	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`

	OrderIndex int `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(wo *Workout) (*Workout, error)
	UpdateWorkout(wo *Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
	DeleteWorkoutById(id int64) error
}

func (pgStore *PostgresWorkoutStore) CreateWorkout(wo *Workout) (*Workout, error) {
	tx, err := pgStore.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("CreateWorkout: %w", err)
	}
	defer tx.Rollback()

	query :=
		`INSERT INTO workouts (title, description, duration_minutes, calories_burned)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	err = tx.QueryRow(query, wo.Title, wo.Description, wo.DurationMinutes, wo.CaloriesBurned).Scan(&wo.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateWorkout Workout: %w", err)
	}

	for _, v := range wo.Entries {
		query := `
			INSERT INTO workout_entries 
				(workout_id, name, notes, sets, reps, duration_seconds, weight, order_index)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err := tx.QueryRow(
			query,
			wo.ID,
			v.Name,
			v.Notes,
			v.Sets,
			v.Reps,
			v.DurationSeconds,
			v.Weight,
			v.OrderIndex,
		).Scan(&v.ID)

		if err != nil {
			return nil, fmt.Errorf("CreateWorkout Entry: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("CreateWorkout Commit: %w", err)
	}

	return wo, nil
}

func (pgStore *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	query := "SELECT id, title, description, duration_minutes, calories_burned FROM workouts WHERE id=$1"

	wo := Workout{}
	err := pgStore.db.QueryRow(query, id).Scan(&wo.ID, &wo.Title, &wo.Description, &wo.DurationMinutes, &wo.CaloriesBurned)
	if err != nil {
		return nil, fmt.Errorf("GetWorkoutById QueryRow: %w", err)
	}

	query =
		"SELECT id, name, notes, sets, reps, weight, duration_seconds, order_index FROM workout_entries WHERE workout_id=$1"

	result, err := pgStore.db.Query(query, wo.ID)
	if err != nil {
		return nil, fmt.Errorf("GetWorkoutById Query: %w", err)
	}
	defer result.Close()

	for result.Next() {
		entry := WorkoutEntry{}
		entry.WorkoutId = wo.ID
		err := result.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Notes,
			&entry.Sets,
			&entry.Reps,
			&entry.Weight,
			&entry.DurationSeconds,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, fmt.Errorf("GetWorkoutById Scan: %w", err)
		}

		wo.Entries = append(wo.Entries, entry)
	}

	return &wo, nil
}

func (pgStore *PostgresWorkoutStore) UpdateWorkout(wo *Workout) (*Workout, error) {
	tx, err := pgStore.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("UpdateWorkout Begin: %w", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	WHERE id = $5
	`
	if _, err := tx.Exec(query, wo.Title, wo.Description, wo.DurationMinutes, wo.CaloriesBurned, wo.ID); err != nil {
		return nil, fmt.Errorf("UpdateWorkout Exec: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM workout_entries WHERE workout_id = $1", wo.ID); err != nil {
		return nil, fmt.Errorf("UpdateWorkout Delete Entry Exec: %w", err)
	}

	for i, v := range wo.Entries {
		query := `
			INSERT INTO workout_entries 
				(workout_id, name, notes, sets, reps, duration_seconds, weight, order_index)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, workout_id
		`
		err := tx.QueryRow(
			query,
			wo.ID,
			v.Name,
			v.Notes,
			v.Sets,
			v.Reps,
			v.DurationSeconds,
			v.Weight,
			v.OrderIndex,
		).Scan(&wo.Entries[i].ID, &wo.Entries[i].WorkoutId)

		if err != nil {
			return nil, fmt.Errorf("UpdateWorkout Entry Creation: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("UpdateWorkout Commit: %w", err)
	}

	return wo, nil
}

func (pgStore *PostgresWorkoutStore) DeleteWorkoutById(id int64) error {
	_, err := pgStore.db.Exec("DELETE FROM workouts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("DeleteWorkoutById Exec: %w", err)
	}
	return nil
}
