package singleton

import "sync"

type Singleton struct {
	instance    sync.Once
}

func (s *Singleton) Construct(init func())  {
	s.instance.Do(init)
}