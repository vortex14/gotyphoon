package awaitable

import "sync"

type Object struct {
	await sync.WaitGroup
}

func (o *Object) Init()  {
	o.await = sync.WaitGroup{}
}

func (o *Object) Add()  {
	o.await.Add(1)
}

func (o *Object) Done()  {
	o.await.Done()
}

func (o *Object) Await()  {
	o.await.Wait()
}