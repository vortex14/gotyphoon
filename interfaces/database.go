package interfaces

import (
	"bufio"
	"os"
)

type Database interface {
	Import(Database string, collection string, inputFile string) (error, uint64)
	Export(Database string, collection string, outFile string) (*bufio.Writer, *os.File, int64, error)
}
