package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
)

//go:generate go tool github.com/luca-arch/asynq-codegen .

const (
	redisDSN           = "redis:6379"
	workersConcurrency = 2
)

// asynq:task task_1
type Task01 struct {
	Created time.Time
}

// asynq:task task_2
type Task02 struct {
	Created time.Time
}

func Consumer(ctx context.Context) error {
	processors := &Processors{
		// Process Task01
		Task01: func(ctx context.Context, task01 *Task01, headers map[string]string) error {
			slog.InfoContext(ctx, "processing task", "created", task01.Created.String(), "headers", headers, "type", TypeTask01)

			return nil
		},
		// Process Task02
		Task02: func(ctx context.Context, task02 *Task02, headers map[string]string) error {
			slog.InfoContext(ctx, "processing task", "created", task02.Created.String(), "headers", headers, "type", TypeTask02)

			return nil
		},
	}

	// Never pass [context.Background] to [Processors.Run]!
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	slog.InfoContext(ctx, "starting Consumer")

	return processors.Run(
		ctx,
		asynq.RedisClientOpt{Addr: redisDSN},
		asynq.Config{Concurrency: workersConcurrency},
	)
}

func Producer(ctx context.Context) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisDSN})
	defer client.Close()

	headers := map[string]string{
		"start": time.Now().String(),
	}

	dispatcher := &Dispatcher{Client: client}

	for {
		// Enqueue Task01 (with context.Background(), no headers)
		task1, info, err := dispatcher.EnqueueTask01(&Task01{
			Created: time.Now(),
		})
		if err != nil {
			return err
		}

		slog.Info("new task pushed", "task.id", info.ID, "task.type", task1.Type(), "type", TypeTask01)

		<-time.After(time.Second)

		// Enqueue Task02 (pass context and headers)
		task2, info, err := dispatcher.EnqueueTask02ContextWithHeaders(ctx, &Task02{
			Created: time.Now(),
		}, headers)
		if err != nil {
			return err
		}

		slog.Info("new task pushed", "task.id", info.ID, "task.type", task2.Type(), "type", TypeTask02)

		<-time.After(time.Second)
	}
}

func main() {
	var err error

	switch os.Args[1] {
	case "consumer":
		err = Consumer(context.Background())
	case "producer":
		err = Producer(context.Background())
	}

	if err != nil {
		panic(err)
	}
}
