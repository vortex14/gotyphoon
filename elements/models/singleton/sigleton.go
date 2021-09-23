package singleton

import "sync"

type Singleton struct {
	instance    sync.Once
	exitOnce    sync.Once
}

func (s *Singleton) Construct(init func())  {
	s.instance.Do(init)
}

func (s *Singleton) Destruct(finish func())  {
	s.exitOnce.Do(finish)
}