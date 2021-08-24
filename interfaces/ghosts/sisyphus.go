package ghosts

import "sync"

type SisyphusInterface interface {
	Roll()
	Stop()
	RollOut()
}

type Sisyphus struct {
	promise sync.WaitGroup
}

func (s *Sisyphus) AddPromise()  {
	s.promise.Add(1)
}

func (s *Sisyphus) PromiseDone()  {
	s.promise.Done()
}

func (s *Sisyphus) WaitPromises()  {
	s.promise.Wait()
}