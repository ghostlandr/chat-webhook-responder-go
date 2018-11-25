package timer

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	defaultLabel = "default"

	ErrAlreadyStarted = errors.New("you've already started that timer")
	ErrAlreadyEnded   = errors.New("you've already ended that timer")
	ErrEndNoStart     = errors.New("you should start that timer before you end it")
	ErrNotFinished    = errors.New("you should run that timer before getting elapsed time 👍")
)

type timer struct {
	starts map[string]time.Time
	ends   map[string]time.Time
}

func (t *timer) Start(label string) error {
	if label == "" {
		label = defaultLabel
	}
	if !t.starts[label].IsZero() {
		return ErrAlreadyStarted
	}
	t.starts[label] = time.Now()
	return nil
}

func (t *timer) End(label string) error {
	if label == "" {
		label = defaultLabel
	}
	if !t.ends[label].IsZero() {
		return ErrAlreadyEnded
	}
	if t.starts[label].IsZero() {
		return ErrEndNoStart
	}
	t.ends[label] = time.Now()
	return nil
}

func (t *timer) Elapsed(label string) (time.Duration, error) {
	if label == "" {
		label = defaultLabel
	}
	if t.starts[label].IsZero() || t.ends[label].IsZero() {
		return 0, ErrNotFinished
	}
	return t.ends[label].Sub(t.starts[label]) / time.Millisecond, nil
}

func (t *timer) String() string {
	bigO := "Timers:\n"
	ct := make(chan string)
	quit := make(chan int)
	go t.formatStrings(ct, quit)
	for {
		select {
		case s := <-ct:
			bigO += s + "\n"
		case <-quit:
			return bigO
		}
	}

	return bigO
}

func (t *timer) formatStrings(ct chan string, quit chan int) {
	var wg sync.WaitGroup
	for label := range t.starts {
		label := label
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			o := fmt.Sprintf("%s: ", label)
			if !t.ends[label].IsZero() {
				el, _ := t.Elapsed(label)
				o += fmt.Sprintf("%dms", el)
			}
			if t.ends[label].IsZero() {
				o += fmt.Sprint("Still running")
			}
			ct <- o
		}()
	}
	wg.Wait()
	quit <- 1
}

func New() *timer {
	return &timer{
		starts: make(map[string]time.Time),
		ends:   make(map[string]time.Time),
	}
}
