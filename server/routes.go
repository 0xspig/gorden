package server

import (
	"net/http"
)

func (app *Application) routes() *http.ServeMux {

	app.logger.Debug("running routes")

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /graph-json", app.getJSON)
	mux.HandleFunc("GET /node-data/{id}", app.getPostHTML)
	mux.HandleFunc("GET /node-links/{id}", app.getLinksHTML)

	return mux
}
