package multiconfig

import "github.com/fatih/structs"

var (
	// DefaultDefaultTag is the default tag name for struct fields to define
	// default values for a field. Example:
	//
	//   // Field's default value is "koding".
	//   Name string `default:"koding"`
	//
	DefaultDefaultTag = "default"
)

type TagLoader struct {
}

func (t *TagLoader) Load(s interface{}) error {
	for _, field := range structs.Fields(s) {
		defaultVal := field.Tag(DefaultDefaultTag)
		if defaultVal == "" {
			continue
		}

		err := fieldSet(field, defaultVal)
		if err != nil {
			return err
		}
	}

	return nil
}
