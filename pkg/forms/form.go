package forms

import "net/url"

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
