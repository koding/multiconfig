package multiconfig

import (
	"strings"
	"testing"
)

func TestValidators(t *testing.T) {
	s := getDefaultServer()
	s.Name = ""

	err := (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("Name should be required")
	}
}

func TestValidatorsEmbededStruct(t *testing.T) {
	s := getDefaultServer()
	s.Postgres.Port = 0

	err := (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("Port should be required")
	}
}

func TestValidatorsCustomTag(t *testing.T) {
	s := getDefaultServer()

	validator := (&RequiredValidator{
		TagName:  "customRequired",
		TagValue: "yes",
	})

	// test happy path
	err := validator.Validate(s)
	if err != nil {
		t.Fatal(err)
	}

	// validate sad case
	s.Postgres.Port = 0
	err = validator.Validate(s)
	if err == nil {
		t.Fatal("Port should be required")
	}

	errStrPrefix := "1 error occurred:"
	errStrSufix := "field 'Postgres.Port' is required"

	if !strings.HasPrefix(err.Error(), errStrPrefix) {
		t.Fatalf("Err string is wrong: expected prefix %s, got: %s", errStrPrefix, err.Error())
	}

	if !strings.HasSuffix(err.Error(), errStrSufix) {
		t.Fatalf("Err string is wrong: expected suffix %s, got: %s", errStrSufix, err.Error())
	}
}
