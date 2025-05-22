package application

import "database/sql"

func store(db *sql.DB, app *Application) error {
	return CreateApplication(db, app)
}

func show(db *sql.DB, id int) (*Application, error) {
	return GetApplicationByID(db, id)
}