package interfaces

type TaskInterface interface {
	IsMaxRetry() bool
	UpdateRetriesCounter()
	IsRetry()bool
}