// Filename: ./cmd/api/healthcheck.go
package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//js := `{"status": "available", "environment" :%q, "version": %q}`
	//js = fmt.Sprintf(js, app.config.env, version)

	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.WriteJSON(w, http.StatusOK, data, nil) //send ok to say everything is okay. We send w caz it sends stuff to the browser.
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
