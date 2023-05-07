// Filename: cmd/api/inputs.go
package main

import (
	"fmt"
	"net/http"

	"github.com/abelwhite/4191/internal/data"
	"github.com/abelwhite/4191/internal/validator"
)

func (app *application) createCourseHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a input that will be provided to us
	// via the request
	var input struct {
		CourseCode   string `json:"name"`
		CourseTitle  string `json:"level"`
		CourseCredit string `json:"contact"`
	}
	// decode our the JSON request
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to the new School struct
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

// func (app *application) showSchoolHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParams(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	//fetch the school with the associated idea id
// 	school, err := app.models.Schools.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"school": school}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// }

// func (app *application) updateSchoolHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParams(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	//fetch the original School
// 	school, err := app.models.Schools.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	var input struct {
// 		Name    string   `json:"name"`
// 		Level   string   `json:"level"`
// 		Contact string   `json:"contact"`
// 		Phone   string   `json:"phone"`
// 		Email   string   `json:"email"`
// 		Website string   `json:"website"`
// 		Address string   `json:"address"`
// 		Mode    []string `json:"mode"`
// 	}
// 	// decode our the JSON request
// 	err = app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	} //Update the original school
// 	//with the new school
// 	school.Name = input.Name
// 	school.Level = input.Level
// 	school.Contact = input.Contact
// 	school.Phone = input.Phone
// 	school.Email = input.Email
// 	school.Website = input.Website
// 	school.Address = input.Address
// 	school.Mode = input.Mode
// 	// let's validate our JSON input
// 	v := validator.New()
// 	// Check for validation errrors
// 	if data.ValidateSchool(v, school); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	//perform the update
// 	err = app.models.Schools.Update(school)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	} //write the response
// 	err = app.writeJSON(w, http.StatusOK, envelope{"school": school}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
