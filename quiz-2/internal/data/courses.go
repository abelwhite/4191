// Filename: ./internal/data/courses.go
package data

import (
	"time"
)

// school represents one row of data in our schools table
type Course struct { //we can get data from client and put it in here and send to db or vise versa
	ID           int64     `json:"id"`
	CourseCode   string    `json:"course_code"`
	CourseTitle  string    `json:"course_title"`
	CourseCredit string    `json:"course_credit"`
	CreatedAt    time.Time `json:"-"`
	Version      int32     `json:"version"`
}
