package interfaces

type CloudManagement interface {
	Deploy()
}

type TestData interface {
	GetFields()
}

type goPromise interface {
	AddPromise()
	PromiseDone()
	WaitPromises()
}








