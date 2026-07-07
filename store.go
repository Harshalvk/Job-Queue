package jobqueue

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) RecordCreated(ctx context.Context, job *Job) error {
	_, error := s.db.Exec(ctx, `
		 INSERT INTO job_history (id, type, payload, status, attempts, max_attempts, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (id) DO NOTHING
	`, job.ID, job.Type, job.Payload, job.Status, job.Attempts, job.MaxAttempts, job.CreatedAt)

	if error != nil {
		return fmt.Errorf("record created %w", error)
	}

	return  nil
}

func (s *Store) RecordStatus (ctx context.Context, job *Job) error {
	_, error := s.db.Exec(ctx, `
		 UPDATE job_history
		 SET status = $2, attempts = $3, last_error = $4, updated_at = now()
		 WHERE id = $1
	`, job.ID, job.Status, job.Attempts, job.LastError)

	if error != nil {
		return fmt.Errorf("record status: %w", error)
	}

	return  nil
}