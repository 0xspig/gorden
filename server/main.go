package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"go.tsmckee.com/cmd/models"
)

type Application struct {
	logger *slog.Logger
	Garden *models.Garden
}

func main() {
	g := models.CreateGarden()

	// TODO make content dir in config or something to search files in
	// for now Im just going to hack in static
	g.PopulateGardenFromDir("ui/content")
	g.ParseAllConnections()
	g.GenAssets()

	addr := flag.String("addr", "localhost:3000", "HTTP Network Address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &Application{
		logger: logger,
		Garden: g,
	}
	logger.Info("Starting server on", "addr", *addr)
	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
