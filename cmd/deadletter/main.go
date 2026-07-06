package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/harshalvk/jobqueue"
	"github.com/redis/go-redis/v9"
)

func main() {
	 action := flag.String("action", "list", "list | requeue | purge")
	 jobID := flag.String("id", "", "job ID (required for requeue)")
	 flag.Parse()

	 rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	 q := jobqueue.NewQueue(rdb)
	 ctx := context.Background()

	 switch *action {
	case "list":
		jobs, error := q.ListDeadLetter(ctx, 50)
		if error != nil {
			log.Fatal(error)
		}
		for _, job := range jobs {
			fmt.Printf(("id=%s type=%s attempts=%d error=%q\n"), job.ID, job.Type, job.Attempts, job.LastError)
		}
	case "requeue":
		if *jobID == "" {
			log.Fatal("--id required for requeue")
		}
		if error := q.RequeueDeadLetter(ctx, *jobID); error != nil {
			log.Fatal(error)
		}
		fmt.Println("requeued: ", *jobID)
	case "purge":
		if error := q.PurgeDeadLetter(ctx); error != nil {
			log.Fatal(error)
		}
		fmt.Println("dead letter queue purged")
	 }
}