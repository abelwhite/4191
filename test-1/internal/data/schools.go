package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/abelwhite/4191/internal/validator"
	"github.com/lib/pq"
)

// School represents one row of data in our schools table
type School struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Level    string    `json:"level"`
	Contact  string    `json:"contact"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Website  string    `json:"website"`
	Address  string    `json:"address"`
	Mode     []string  `json:"mode"`
	CreateAt time.Time `json:"-"`
	Version  int32     `json:"version"`
}

func ValidateSchool(v *validator.Validator, school *School) {
	// Use the Check() method to execute our validation checks
	v.Check(school.Name != "", "name", "must be provided")
	v.Check(len(school.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(school.Level != "", "level", "must be provided")
	v.Check(len(school.Level) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(school.Contact != "", "contact", "must be provided")
	v.Check(len(school.Contact) <= 200, "contact", "must not be more than 200 bytes long")

	v.Check(school.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(school.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(school.Email != "", "email", "must be provided")
	v.Check(validator.Matches(school.Email, validator.EmailRX), "email", "must be a valid email address")

	v.Check(school.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(school.Website), "website", "must be a valid URL")

	v.Check(school.Address != "", "address", "must be provided")
	v.Check(len(school.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(school.Mode != nil, "mode", "must be provided")
	v.Check(len(school.Mode) >= 1, "mode", "must contain at least 1 entry")
	v.Check(len(school.Mode) <= 5, "mode", "must contain at most 5 entries")
	v.Check(validator.Unique(school.Mode), "mode", "must not contain duplicate entries")
}

// implement our models
// 1. Define our model
type SchoolModel struct {
	DB *sql.DB
}

// insert a new school
func (m SchoolModel) Insert(school *School) error {
	//Write an sql quote to insert
	query := `
		INSERT INTO schools(name, level, contact, phone, email, website, address, mode)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version
	`
	//collect the data fields into a slice
	args := []interface{}{school.Name, school.Level, school.Contact, school.Phone, school.Email, school.Website, school.Address, school.Mode}

	return m.DB.QueryRow(query, args...).Scan(&school.ID, &school.CreateAt, &school.Version)

}

// Get a school
func (m SchoolModel) Get(id int64) (*School, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, name, level, contact, phone, email, website, address, mode, version
		FROM schools
		WHERE id = $1
	`
	//define a School variable to hold the school return
	var school School
	err := m.DB.QueryRow(query, id).Scan(
		&school.ID,
		&school.CreateAt,
		&school.Name,
		&school.Level,
		&school.Contact,
		&school.Phone,
		&school.Email,
		&school.Website,
		&school.Address,
		pq.Array(&school.Mode),
		&school.Version,
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
	return &school, nil

}

// update school data
func (m SchoolModel) Update(school *School) error {
	//create
	query := `
		UPDATE schools 
		SET name = $1, level = $2, content = $3, phone = $4, email = $5, website = $6, 
		address = $7, mode = $8, version = version + 1
		WHERE id = $9
		RETURNING version 
		`
	args := []interface{}{
		school.Name,
		school.Level,
		school.Contact,
		school.Phone,
		school.Email,
		school.Website,
		school.Address,
		pq.Array(school.Mode),
		school.ID,
	}
	return m.DB.QueryRow(query, args...).Scan(&school.Version)
}

// delete school data
func (m SchoolModel) Delete(id int64) error {
	return nil
}
