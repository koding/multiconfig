package multiconfig

// Validator validates the config against any predefined rules, those predefined
// rules should be given to this package. The implementer will be responsible
// about the logic
type Validator interface {
	// Validate validates the config struct
	Validate(s interface{}) error
}

// DefaultValidator implements the Validator interface
type DefaultValidator struct {
	Validators []Validator
}

func (d *DefaultValidator) Validate(s interface{}) error {
	if len(d.Validators) == 0 {
		return nil
	}

	return nil
}

func NewValidator(validators ...Validator) *DefaultValidator {
	return &DefaultValidator{
		Validators: validators,
	}
}
