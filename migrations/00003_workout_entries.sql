
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER,
    duration_seconds INTEGER,
    weight DECIMAL(5,2),
    notes TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_workout_entry CHECK (
        (duration_seconds IS NOT NULL OR reps IS NOT NULL) AND 
        (duration_seconds IS NULL OR reps IS NULL)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_entries;
-- +goose StatementEnd