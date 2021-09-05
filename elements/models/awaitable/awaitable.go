package awaitable

import "sync"

type Object struct {
	await sync.WaitGroup
}

func (o *Object) Await()  {
	o.await.Wait()
}