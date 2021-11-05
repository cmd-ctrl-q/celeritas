package celeritas

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (c *Celeritas) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

// Valid validates a validation.
// When there are no errors, return true, else true
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error to the validation type
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Has checks if a field is in the request
func (v *Validation) Has(field string, r *http.Request) (has bool) {
	x := r.Form.Get(field)

	if x != "" {
		has = true
	}

	return
}

// Required checks for required fields
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		// check value in field
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "This field cannot be blank")
		}
	}
}

// Check checks that form posts and posted json on api are ok
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// IsEmail checks for valid email address
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

// IsInt checks that a value on a form post is an integer
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be an integer")
	}
}

// IsFloat checks that a value on a form post is a float
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

// IsDateISO checks for the correct date time format
func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.AddError(field, "This field must be a date in the form of YYYY-MM-DD")
	}
}

// NoSpaces checks for white spaces in a value
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "Spaces are not permitted")
	}
}
