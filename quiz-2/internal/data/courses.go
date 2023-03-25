// Filename: ./internal/data/courses.go
package data

import (
	"time"
)

// school represents one row of data in our schools table
type Course struct { //we can get data from client and put it in here and send to db or vise versa
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}
