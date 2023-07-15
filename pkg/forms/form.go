package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// New Init a Form struct
func New(data url.Values) *Form {
	return &Form{
		data, errors(map[string][]string{}),
	}
}

// Required Implement a required method to check that specific fields in the data form are present and not blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

// MaxLength Implement a MaxLength()
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)

	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters", d))
	}
}

// PermittedValues Implement PermittedValues method to check that a specific field in the form matches appropriate value
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Implement a valid() to return if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
