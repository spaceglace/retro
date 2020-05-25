package main

import (
	"retro/config"
	"retro/retro"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	logger, _ = zap.NewDevelopment()

	settings := config.NewConfig()
	retro := retro.NewRetro()

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     800,
		Height:    600,
		Resizable: true,
		Title:     "Retro",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
	})

	app.Bind(settings)
	app.Bind(retro)
	app.Run()
}
