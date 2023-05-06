package cmd

import (
	"go-pomodoro/pomodoro"
	"go-pomodoro/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
