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

// Validate tries to validate given struct with all the validators. If it doesnt
// have any Validator it will simply skip the validation step. If any of the
// given validators return err, it will stop validating and return it.
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

// MustValidate validates the struct, it panics if gets any error
func (d *DefaultValidator) MustValidate(s interface{}) {
	if err := d.Validate(s); err != nil {
		panic(err)
	}
}

// NewValidator accepts variadic validators and satisfies Validator interface.
func NewValidator(validators ...Validator) *DefaultValidator {
	return &DefaultValidator{
		Validators: validators,
	}
}

// RequiredValidator validates the struct against zero values
type RequiredValidator struct {
	//  TagName holds the validator tag name
	TagName string

	// TagValue holds the expected value of the validator
	TagValue string
}

// Validate validates the given struct agaist field's zero values. If
// intentionaly, the value of a field is `zero-valued`(e.g false, 0, "")
// required tag should not be set for that field.
func (e *RequiredValidator) Validate(s interface{}) error {
	if e.TagName == "" {
		e.TagName = "required"
	}

	if e.TagValue == "" {
		e.TagValue = "true"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *RequiredValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		val := field.Tag(e.TagName)
		if val != e.TagValue {
			return nil
		}

		if field.IsZero() {
			return fmt.Errorf("field %s is required", fieldName)
		}
	}

	return nil
}
