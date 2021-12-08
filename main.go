package main

import (
	"embed"
	"log"
	"os"
)

//go:embed static
var staticFS embed.FS

//go:embed static/img/logo.webp
var logo []byte

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	logger.Println("Hello world")
}
