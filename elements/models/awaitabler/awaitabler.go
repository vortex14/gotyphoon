package awaitabler

import "sync"

type Object struct {
	awaitable sync.WaitGroup
}

func (o *Object) Init()  {
	o.awaitable = sync.WaitGroup{}
}

func (o *Object) Add()  {
	o.awaitable.Add(1)
}

func (o *Object) Done()  {
	o.awaitable.Done()
}

func (o *Object) Await()  {
	o.awaitable.Wait()
}