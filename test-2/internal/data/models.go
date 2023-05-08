// filename: internal/data/models.go
package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Courses CourseModel
}

// this function creates a models instance
func NewModels(db *sql.DB) Models {
	return Models{
		Courses: CourseModel{DB: db},
	}
}
