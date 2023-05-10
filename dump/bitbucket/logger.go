package bitbucket

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "bitbucket | ", log.Ldate|log.Ltime|log.Lmicroseconds)
)
