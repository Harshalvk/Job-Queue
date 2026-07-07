package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harshalvk/jobqueue"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	q := jobqueue.NewQueue(rdb)
	ctx := context.Background()

	db, error := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if error != nil {
		panic(error)
	}
	defer db.Close()
	store := jobqueue.NewStore(db)

	payload, _ := json.Marshal(map[string]string{"to": "devwork2004@gmail.com"})
	job := jobqueue.NewJob("send_email", payload, 3)

	if err := q.Enqueue(ctx, job); err != nil {
		panic(err)
	}
	if error := store.RecordCreated(ctx, job); error != nil {
		panic(error)
	}
	fmt.Println("enqueued:", job.ID)
}