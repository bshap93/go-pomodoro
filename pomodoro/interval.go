package pomodoro

import (
	"context"
	"errors"
	"time"
)

const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

// / State constants
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ByID(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int) ([]Interval, error)
}

// / As an exercise for later, this list of errors is incomplete.
var (
	ErrNoIntervals         = errors.New("No intervals")
	ErrIntervalsNotRunning = errors.New("Intervals not running")
	ErrIntervalsCompleted  = errors.New("Interval is completed or cancelled")
	ErrInvalidState        = errors.New("Invalid State")
	ErrInvalidID           = errors.New("Invalid ID")
)

type IntervalConfig struct {
	repo               Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

func NewConfig(repo Repository, pomodoro, shortBreak,
	longBreak time.Duration) *IntervalConfig {

	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}

	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}

	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}

	return c
}

func nextCategory(r Repository) (string, error) {
	// last interval just before the current
	li, err := r.Last()
	if err != nil && err == ErrNoIntervals {
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}

	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}

	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}

	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}

	return CategoryLongBreak, nil
}

type Callback func(Interval)

func tick(ctx context.Context, id int64, config *IntervalConfig,
  start, periodic, end) error {

  ticker := time.NewTicker(time.Second) 
  defer ticker.Stop()

  i, err := config.repo.ByID(id)
  if err != nil {
    return err
  } 

  if i.State == StatePaused {
    return nil
  }

  expire := time.After(i.PlannedDuration - i.ActualDuration)
  start(i)
  for {
    select {
    // Ticker's channel
    case <-ticker.C:
      i, err := config.repo.ByID(id)
      if err != nil {
        return err
      }
      if i.State == StatePaused {
         return nil
      }
      i.ActualDuration += time.Second
      if err := config.repo.Update(i); err != nil {
         
      }
    }
  }
}
