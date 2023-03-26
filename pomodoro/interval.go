package pomodoro

import (
	"context"
	"errors"
	"time"
)

// Categories of Time Block, the three basic ones
// in any Pomodoro implementation being the work block,
// and the short and long break.
const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

// State constansts.
// Initial -> StateNotStarted
// StateDone -> Means things are stopped, as does StateNotStarted,
// StatePaused, and StateCancelled.
// StateRunning -> Timer is ticking.
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Interval or Time Block. Each has a unique ID.
// There is a amount of time each interval has when it starts.
// There is a planned duration based on what state the interval
// is in. The ActualDuration is retroactive, and is only set when
// state is StateDone, StateCancelled.
// Contains a reference to a category that an interval is in,
// and the state it's in at a given moment.
type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

// Interface for Create to create an interval and returnits ID.
// Update to update the interval.
// ByID to retrieve an interval by passing in its id.
// Last to find the last interval.
// Breaks to retrieve intervals that are breaks.
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
	start, periodic, end Callback) error {

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
				return err
			}
			periodic(i)
		case <-expire:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			i.State = StateDone
			end(i)
			return config.repo.Update(i)
		case <-ctx.Done():
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			i.State = StateCancelled
			return config.repo.Update(i)
		}
	}
}

func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return i, err
	}
	i.Category = category
	switch category {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}
	// We're expecting an error if the interval is already created.
	// We still will return its ID.
	if i.ID, err = config.repo.Create(i); err != nil {
		return i, err
	}
	// In this case, we created an interval and return it.
	return i, nil
}

/// Public API

func GetInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	var err error
	i, err = config.repo.Last()
	if err != nil && err != ErrNoIntervals {
		return i, err
	}

	if err == nil && i.State != StateCancelled && i.State != StateDone {
		return i, err
	}

	return newInterval(config)
}
