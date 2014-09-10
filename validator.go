package multiconfig

import (
	"fmt"

	"github.com/fatih/structs"
)

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

type RequiredValidator struct {
	RequiredTagName string
}

func (e *RequiredValidator) Validate(s interface{}) error {
	if e.RequiredTagName == "" {
		e.RequiredTagName = "required"
	}

	requiredFields := []string{}
	for _, field := range structs.Fields(s) {
		defaultVal := field.Tag(e.RequiredTagName)
		if defaultVal == "" {
			continue
		}

		if field.IsZero() {
			requiredFields = append(requiredFields, field.Name())
		}
	}

	if len(requiredFields) == 0 {
		return nil
	}

	return fmt.Errorf("Field(s) [%v] are required", requiredFields)
}
