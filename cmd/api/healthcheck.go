package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// create map to hold app status info
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil) // helper.go function to write JSON
	if err != nil {
		// something went wrong, throw an error
		app.serverErrorResponse(w, r, err)
	}
}
