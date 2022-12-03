package ssgo

import (
	"log"
	"os"
)

var (
	HackerLogger = log.New(os.Stderr, "[HACKER] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	ServerLogger = log.New(os.Stderr, "[SERVER] ", log.Ldate|log.Ltime|log.Lmsgprefix)
)
