package repository

import (
	"fmt"
	"go-pomodoro/pomodoro"
	"sync"
)

type inMemoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		intervals: []pomodoro.Interval{},
	}
}
