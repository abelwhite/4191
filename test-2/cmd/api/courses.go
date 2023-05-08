// Filename: cmd/api/inputs.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abelwhite/4191/internal/data"
	"github.com/abelwhite/4191/internal/validator"
)

func (app *application) createCourseHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a input that will be provided to us
	// via the request
	var input struct {
		CourseCode   string `json:"course_code"`
		CourseTitle  string `json:"course_title"`
		CourseCredit string `json:"course_Credit"`
	}
	// decode our the JSON request
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to the new course struct
	course := &data.Course{
		CourseCode:   input.CourseCode,
		CourseTitle:  input.CourseTitle,
		CourseCredit: input.CourseCredit,
	}
	// let's validate our JSON input
	v := validator.New()
	// Check for validation errrors
	if data.ValidateCourse(v, course); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//write our validated school to database
	err = app.models.Courses.Insert(course)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//set the creation header
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/course/%d", course.ID))
	//write the response

	err = app.writeJSON(w, http.StatusCreated, envelope{"course": course}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} //write the response

	// err = app.writeJSON(w, http.StatusCreated, envelope{" course": course}, headers)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }
}

func (app *application) showCourseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//fetch the school with the associated idea id
	course, err := app.models.Courses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"course": course}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateCourseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//fetch the original School
	course, err := app.models.Courses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//create an input to hold the read in from fetch client
	//we will not use string but pointers for partial updates
	//pointers have a default value of nil
	var input struct {
		CourseCode   *string `json:"course_code"`
		CourseTitle  *string `json:"course_title"`
		CourseCredit *string `json:"course_credit"`
	}
	// decode our the JSON request
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//check for updates
	if input.CourseCode != nil {
		course.CourseCode = *input.CourseCode
	}
	if input.CourseTitle != nil {
		course.CourseTitle = *input.CourseTitle
	}
	if input.CourseCredit != nil {
		course.CourseCredit = *input.CourseCredit
	}

	// let's validate our JSON input to make sure they are valid
	v := validator.New()
	// Check for validation errrors
	if data.ValidateCourse(v, course); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	//perform the update
	err = app.models.Courses.Update(course)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	} //write the response
	err = app.writeJSON(w, http.StatusOK, envelope{"course": course}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//delete the schools from the database
	err = app.models.Courses.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//if no errors then deletion was succesfull
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "courses was deleted successful"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
