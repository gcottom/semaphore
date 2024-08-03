package semaphore

import "sync"

type Semaphore struct {
	Channel   chan struct{}
	WaitGroup *sync.WaitGroup
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		Channel:   make(chan struct{}, n),
		WaitGroup: &sync.WaitGroup{},
	}
}

func (s *Semaphore) Acquire() {
	s.Channel <- struct{}{}
	s.WaitGroup.Add(1)
}

func (s *Semaphore) Release() {
	<-s.Channel
	s.WaitGroup.Done()
}

func (s *Semaphore) Wait() {
	s.WaitGroup.Wait()
}
