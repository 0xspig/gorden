package server

import (
	"log/slog"
	"net/http"
	"os"

	"gorden.tsmckee.com/garden"
)

type Application struct {
	logger *slog.Logger
	Garden *garden.Garden
}

func (app *Application) Init(dir string, renderDrafts bool) {

	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	os.Chdir(dir)

	g := garden.CreateGarden()

	// TODO make content dir in config or something to search files in
	// for now Im just going to hack in static
	g.PopulateGardenFromDir("ui/content")
	g.ParseAllConnections()
	g.GenAssets()

	newLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	*app = Application{
		logger: newLogger,
		Garden: g,
	}
}

func (app *Application) Start(addr string) {
	app.logger.Info("Starting server on", "addr", addr)
	err := http.ListenAndServe(addr, app.routes())

	app.logger.Error(err.Error())
	os.Exit(1)
}
