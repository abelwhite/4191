package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/abelwhite/4191/internal/validator"
)

// Course represents one row of data in our Courses table
type Course struct {
	ID           int64     `json:"id"`
	CourseCode   string    `json:"course_code"`
	CourseTitle  string    `json:"course_title"`
	CourseCredit string    `json:"course_credit"`
	CreatedAt    time.Time `json:"-"`
	Version      int32     `json:"version"`
}

func ValidateCourse(v *validator.Validator, course *Course) {
	// Use the Check() method to execute our validation checks
	v.Check(course.CourseCode != "", "course_code", "must be provided")
	v.Check(len(course.CourseCode) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(course.CourseTitle != "", "course_title", "must be provided")
	v.Check(len(course.CourseTitle) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(course.CourseCredit != "", "course_credit", "must be provided")
	v.Check(len(course.CourseCredit) <= 200, "contact", "must not be more than 200 bytes long")
}

// implement our models
// 1. Define our model
type CourseModel struct {
	DB *sql.DB
}

// insert a new course
func (m CourseModel) Insert(course *Course) error {
	//Write an sql quote to insert
	query := `
		INSERT INTO courses(course_code, course_title, course_credit)
		VALUES($1, $2, $3)
		RETURNING id, created_at, version
	`
	//collect the data fields into a slice
	args := []interface{}{
		course.CourseCode,
		course.CourseTitle,
		course.CourseCredit}

	return m.DB.QueryRow(query, args...).Scan(&course.ID, &course.CreatedAt, &course.Version)

}

// Get a school
func (m CourseModel) Get(id int64) (*Course, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, course_code, course_title, course_credit, version
		FROM courses
		WHERE id = $1
	`
	//define a School variable to hold the school return
	var course Course
	err := m.DB.QueryRow(query, id).Scan(
		&course.ID,
		&course.CreatedAt,
		&course.CourseCode,
		&course.CourseTitle,
		&course.CourseCredit,
		&course.Version,
	)

	//check for errors
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	//no errors found
	return &course, nil

}

// update school data with fixed DATA RACES
// optimistc locking in the version number
func (m CourseModel) Update(course *Course) error {
	//create
	query := `
		UPDATE courses 
		SET course_code = $1, course_title = $2, course_credit = $3, version = version + 1
		WHERE id = $4
		AND version =$5
		RETURNING version 
		`
	args := []interface{}{
		course.CourseCode,
		course.CourseTitle,
		course.CourseCredit,
		course.ID,
		course.Version, //we want to keep a check in version to make sure it hasnt changed
	}
	//check for edit conflicts
	err := m.DB.QueryRow(query, args...).Scan(&course.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows): //tells me no rows was return
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// delete school data
func (m CourseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	//create the delete query
	query := `
	DELETE FROM courses
	WHERE id =$1
	`
	//execute the query
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	//check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	//check if 0 rows affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
