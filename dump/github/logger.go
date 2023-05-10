package github

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "github | ", log.Ldate|log.Ltime|log.Lmicroseconds)
)
