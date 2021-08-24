package interfaces



type Pipeline interface {
	Run()
	Crawl()
	Finish()
	Switch()
	Retry()
	Await()
}

