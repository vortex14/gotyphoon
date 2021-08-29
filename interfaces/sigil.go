package interfaces

import (
	"sync"
)

type SigilOptions struct {
	Name 			string
	Required 		bool
	Description 	string
	AllowedCantrips []string
	CantripsMap 	map[string]CantripInterface
}

type Sigil struct {
	Run 	 bool
	//status 	 bool
	Critical bool
	Options  *SigilOptions
	Conjure  func() error
	promise  sync.WaitGroup
}


type SigilBuilder struct {
	Constructor func(opt *SigilOptions, project Project) SigilInterface
	Options *SigilOptions
}

type SigilInterface interface {
	Crash(err error)
	IsCritical() bool
	GetName() string
	Conjure(project Project, name string) error
	GetCantrips() map[string]CantripInterface
}

func (s *Sigil) AddPromise()  {
	s.promise.Add(1)
}

func (s *Sigil) PromiseDone()  {
	s.promise.Done()
}

func (s *Sigil) WaitPromises()  {
	s.promise.Wait()
}