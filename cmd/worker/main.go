package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/harshalvk/jobqueue"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// func sendEmailHandler(ctx context.Context, job *jobqueue.Job) error {
// 	var payload struct {
// 		To string `json:"to"`
// 	}
// 	if err := json.Unmarshal(job.Payload, &payload); err != nil {
// 		return err
// 	}
// 	fmt.Printf("sending email to %s (job %s)\n", payload.To, job.ID)
// 	return nil
// }

// simulated version to fail a job
func sendEmailHandler(ctx context.Context, job *jobqueue.Job) error {
	return fmt.Errorf("simulated failure")
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	queue := jobqueue.NewQueue(rdb)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, error := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if error != nil {
		panic(error)
	}
	defer db.Close()
	store := jobqueue.NewStore(db)

	pool := jobqueue.NewWorkerPool(queue, store, 5) // 5 concurrent workers
	pool.RegisterHandler("send_email", sendEmailHandler)

	fmt.Println("worker pool started, waiting for jobs...")
	pool.Start(ctx)
	fmt.Println("worker pool stopped")
} 