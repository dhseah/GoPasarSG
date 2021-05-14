package forms

// errors hold the validation error messages for forms.
// The name of the form field will be used as the key.
type errors map[string][]string

// Add appends an error message to the slice of error
// messages for a given field in the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves the first error message from the slice
// of error messages for a given field in the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
