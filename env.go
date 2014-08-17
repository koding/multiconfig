package multiconfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// EnvironmentLoader satisifies the loader interface. It loads the
// configuration from the environment variables in the form of
// STRUCTNAME_FIELDNAME.
type EnvironmentLoader struct{}

func (e *EnvironmentLoader) Load(s interface{}) error {
	strct := structs.New(s)
	strctName := strct.Name()

	for _, field := range strct.Fields() {
		if err := e.processField(strctName, field); err != nil {
			return err
		}
	}

	return nil
}

// processField gets leading name for the env variable and combines the current
// field's name and generates environemnt variable names recursively
func (e *EnvironmentLoader) processField(prefix string, field *structs.Field) error {
	fieldName := generateFieldName(prefix, field)

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		v := os.Getenv(fieldName)
		if v == "" {
			return nil
		}

		if err := fieldSet(field, v); err != nil {
			return err
		}
	}

	return nil
}

// PrintEnvs prints the generated environment variables to the std out.
func (e *EnvironmentLoader) PrintEnvs(s interface{}) {
	strct := structs.New(s)
	strctName := strct.Name()

	for _, field := range strct.Fields() {
		e.printEnvVariables(strctName, field)
	}
}

func (e *EnvironmentLoader) printEnvVariables(prefix string, field *structs.Field) {
	fieldName := generateFieldName(prefix, field)

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			e.printEnvVariables(fieldName, f)
		}
	default:
		fmt.Println("  ", fieldName)
	}
}

func generateFieldName(prefix string, field *structs.Field) string {
	return strings.ToUpper(prefix) + "_" + strings.ToUpper(field.Name())
}
