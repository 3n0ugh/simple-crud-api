package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if Errors map is empty
func (v Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds error if the error is not exists in the Errors map
func (v *Validator) AddError(key, msg string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = msg
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

// In returns true if a specific value is in a list of strings.
func In(value string, list ...string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
