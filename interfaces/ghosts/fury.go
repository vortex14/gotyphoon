package ghosts

import "sync"

type FuryInterface interface {
	
}

type Fury struct {
	promise sync.WaitGroup
}

func (f *Fury) AddPromise()  {
	f.promise.Add(1)
}

func (f *Fury) PromiseDone()  {
	f.promise.Done()
}

func (f *Fury) WaitPromises()  {
	f.promise.Wait()
}
