package main

import (
	"os"

	"github.com/ArthurWang23/miniblog/cmd/mb-userservice/app"
)

func main() {
	if err := app.NewUserServiceCommand().Execute(); err != nil {
		os.Exit(1)
	}
}