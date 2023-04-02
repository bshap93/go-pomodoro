package pomodoro_test

import (
	"testing"

	"go-pomodoro/pomodoro"
	"go-pomodoro/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
