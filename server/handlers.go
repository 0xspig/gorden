package server

import (
	"net/http"
	"os"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("home Function run")
	w.Header().Add("Server", "McServer")

	ts := app.Garden.Templates["home_template"]

	err := ts.ExecuteTemplate(w, "base", app.Garden)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *Application) getJSON(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getJSON Function Run")
	w.Header().Add("Server", "McServer")
	json, err := os.ReadFile("static/gen/graph-data.json")
	if err != nil {
		app.logger.Error("Error loading json", err)
		return
	}

	w.Write(json)
}

func (app *Application) getPostHTML(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getPostHTML function run")
	w.Header().Add("Server", "McServer")
	html := app.Garden.NodeToHTML(r.PathValue("id"))
	w.Write(html)
}

func (app *Application) getLinksHTML(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getPostHTML function run")
	w.Header().Add("Server", "McServer")
	html := app.Garden.NodeLinksToHTML(r.PathValue("id"))
	w.Write(html)
}
