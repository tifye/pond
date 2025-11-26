package main

import (
	"log"
	"os"
	"strings"

	"github.com/tifye/pond/internal/app"
)

type mode = string

const (
	ScreensaverModeFlag mode = "/c"
	ConfigModeFlag      mode = "/p"
	PreviewModeFlag     mode = "/s"
)

func main() {
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch {
		case strings.HasPrefix(arg, ConfigModeFlag):
			return
		case strings.HasPrefix(arg, PreviewModeFlag):
		case strings.HasPrefix(arg, ScreensaverModeFlag):
		default:
		}
	}

	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
