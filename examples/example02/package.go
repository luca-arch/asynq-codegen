package example02

//go:generate go run ../..

// asynq:task    sample_task
// asynq:retry   1
// asynq:timeout 10s
type ExampleTask struct {
	ExampleContent string
}
