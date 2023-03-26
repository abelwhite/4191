// Filename: ./cmd/api/handlers.go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abelwhite/4191/internal/data"
)

func (app *application) createCoursesHandler(w http.ResponseWriter, r *http.Request) {
	//create a struct to hold a school that will be provided to us
	//via the request
	var input struct {
		CourseCode   string `json:"course_code"`
		CourseTitle  string `json:"course_title"`
		CourseCredit string `json:"course_credit"`
	}
	//Decode the JSON request
	err := app.readJSON(w, r, &input) //we take r.Body and decode it into input
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//Print the request
	fmt.Fprintf(w, "%+v\n", input)

}

func (app *application) showCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "School displayed...")
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fmt.Fprintf(w, "Show details of Courses %d \n ", id)
	course := data.Course{
		ID:           id,
		CreatedAt:    time.Now(),
		CourseCode:   "CMPS142",
		CourseTitle:  "Principles of Programming",
		CourseCredit: "3",
		Version:      1,
	}
	err = app.WriteJSON(w, http.StatusOK, envelope{"course": course}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err) //
		return
	}

}
