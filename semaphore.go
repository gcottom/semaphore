package semaphore

import (
	"errors"
	"sync"
)

type Semaphore struct {
	Active    bool
	Channel   chan struct{}
	WaitGroup *sync.WaitGroup
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		Active:    true,
		Channel:   make(chan struct{}, n),
		WaitGroup: &sync.WaitGroup{},
	}
}

func (s *Semaphore) Acquire() error {
	if s.Active {
		s.Channel <- struct{}{}
		s.WaitGroup.Add(1)
		return nil
	}
	return errors.New("Semaphore not active")
}

func (s *Semaphore) Release() {
	if !s.Active {
		return
	}
	<-s.Channel
	s.WaitGroup.Done()
}

func (s *Semaphore) Wait() {
	if !s.Active {
		return
	}
	s.WaitGroup.Wait()
}

func (s *Semaphore) Cancel() {
	if !s.Active {
		return
	}
	s.Active = false
	close(s.Channel)
	drainWaitGroup(s.WaitGroup)
}

func drainWaitGroup(wg *sync.WaitGroup) {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	default:
		wg.Done()
	}
}
