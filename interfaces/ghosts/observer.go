package ghosts

import "sync"

type ObserverInterface interface {
	
}

type Observer struct {
	promise sync.WaitGroup
}

func (o *Observer) AddPromise()  {
	o.promise.Add(1)
}

func (o *Observer) PromiseDone()  {
	o.promise.Done()
}

func (o *Observer) WaitPromises()  {
	o.promise.Wait()
}
