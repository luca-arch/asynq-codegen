package example01_test

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/luca-arch/asynq-codegen/examples/example01"
)

func ExampleNewSendEmailFromJSONBytes() {
	task, err := example01.NewSendEmailFromJSONBytes([]byte(`
		{
			"body": "Hello world",
			"from": "sender@example.com",
			"to":   "recipient@example.com"
		}
	`))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message from %s: %s", task.From, task.Body)

	// output:
	// Message from sender@example.com: Hello world
}

func ExampleEnqueueSendEmailContext() {
	asynqRedisClientOpt := asynq.RedisClientOpt{Addr: "127.0.0.1:6379", DialTimeout: 1}

	client := asynq.NewClient(asynqRedisClientOpt)
	defer client.Close()

	example01.EnqueueSendEmailContext(context.Background(), client, &example01.SendEmail{
		Body: "Hello world",
		From: "sender@example.com",
		To:   "recipient@example.com",
	})

	// output:
	//
}

func ExampleProcessors_HandleAll() {
	serveMux := asynq.NewServeMux()

	if err := (&example01.Processors{
		SendEmail: func(context.Context, *example01.SendEmail) error {
			return nil
		},
		SendSMS: func(context.Context, *example01.SendSMS) error {
			return nil
		},
	}).HandleAll(serveMux); err != nil {
		panic(err)
	}

	// output:
	//
}
