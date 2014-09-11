package multiconfig

import (
	"fmt"
	"reflect"

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

	for _, validator := range d.Validators {
		if err := validator.Validate(s); err != nil {
			return err
		}
	}

	return nil
}

func NewValidator(validators ...Validator) *DefaultValidator {
	return &DefaultValidator{
		Validators: validators,
	}
}

type RequiredValidator struct {
	TagName  string
	TagValue string
}

// Validate validates the given struct agaist field's zero values. By
// If intentionaly, the value of a field is `zero-valued`(e.g false, 0, "")
// required tag should not be set for that field.
func (e *RequiredValidator) Validate(s interface{}) error {
	if e.TagName == "" {
		e.TagName = "required"
	}

	if e.TagValue == "" {
		e.TagValue = "true"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField(field); err != nil {
			return err
		}
	}

	return nil
}

func (e *RequiredValidator) processField(field *structs.Field) error {
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := e.processField(f); err != nil {
				return err
			}
		}
	default:
		val := field.Tag(e.TagName)
		if val != e.TagValue {
			return nil
		}

		if field.IsZero() {
			// todo add parent struct names into error
			return fmt.Errorf("field %s is required", field.Name())
		}
	}

	return nil
}
