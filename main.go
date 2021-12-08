package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	logger.Println("Hello world")
}
