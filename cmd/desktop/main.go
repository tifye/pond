package main

import (
	"log"

	"github.com/tifye/pond/internal/app"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
