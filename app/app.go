package app

import "main.go/db"

// App struct
type App struct {
	Database *db.Database
}

// New app init
func New() (app *App, err error) {
	app = &App{}

	app.Database, err = db.New()
	if err != nil {
		return nil, err
	}

	return app, err
}

// Close db connection
func (a *App) Close() error {
	return a.Database.Close()
}