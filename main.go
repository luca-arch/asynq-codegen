package main

import (
	"flag"
	"log/slog"
	"os"

	asynqcodegen "github.com/luca-arch/async-codegen/asynq-codegen"
)

func main() {
	cwd, _ := os.Getwd()

	fs := flag.NewFlagSet("asynq-codegen", flag.ExitOnError)

	var (
		wd    = fs.String("working-dir", cwd, "working directory")
		quiet = fs.Bool("quiet", false, "disable all output")
	)

	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
	}

	if *quiet {
		slog.SetDefault(slog.New(slog.DiscardHandler))
	}

	if err := asynqcodegen.Main(*wd); err != nil {
		panic(err)
	}
}
