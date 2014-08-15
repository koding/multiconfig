// Package multiconfig provides a way to load and read configurations from
// multiple sources
package multiconfig

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

var (
	errPathNotSet = errors.New("config path is not set")
)

// Config is used to handle multiple configuration sources
type Config struct {
	// Path contains the JSON or TOML file
	Path string
}

// NewWithPath returns a new instance of Config to read from the given
// configuration file.
func NewWithPath(path string) *Config {
	return &Config{
		Path: path,
	}
}

// New returns a new instance of Config.
func New() *Config {
	return &Config{}
}

// MustLoad is like Load but panics if the config cannot be parsed.
func (c *Config) MustLoad(conf interface{}) {
	if err := c.Load(conf); err != nil {
		panic(err)
	}
}

// Load initializes the given pointer of struct s with configuration from the
// default sources. The order of load is LoadFile, LoadEnv and lastly LoadFlag.
// An error in any step stops the loading process. Each step overrides the
// previous step's config (i.e: defining a flag will override previous
// environment or file config). To customize the order use the individual load
// functions.
func (c *Config) Load(conf interface{}) error {
	if !structs.IsStruct(conf) {
		return fmt.Errorf("passed configuration is not a struct: %T", conf)
	}

	// Initialize struct from the config path.
	if err := c.LoadFile(conf); err != nil {
		// Means it its initialzed with New()
		if err != errPathNotSet {
			return err
		}
	}

	// If any environment variable is available override it.
	if err := c.LoadEnv(conf); err != nil {
		return err
	}

	// Finally check if any flag is defined, which overrides the field again
	if err := c.LoadFlag(conf); err != nil {
		return err
	}

	return nil
}

// LoadFile sets the fields of the given s struct by reading from the given path.
func (c *Config) LoadFile(s interface{}) error {
	// Initialize struct from the config path.
	if c.Path == "" {
		return errPathNotSet
	}

	if strings.HasSuffix(c.Path, "toml") {
		if _, err := toml.DecodeFile(c.Path, s); err != nil {
			return err
		}

		return nil
	}

	if strings.HasSuffix(c.Path, "json") {
		file, err := ioutil.ReadFile(c.Path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(file, s); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("file format is not supported: %s. Want TOML or JSON.", c.Path)
}

// LoadEnv sets the fields of the given s struct by looking for environment
// variables in the form of STRUCTNAME_FIELDNAME.
func (c *Config) LoadEnv(s interface{}) error {
	strct := structs.New(s)
	strctName := strct.Name()

	for _, field := range strct.Fields() {
		envName := strings.ToUpper(strctName) + "_" + strings.ToUpper(field.Name())

		v := os.Getenv(envName)
		if v == "" {
			continue
		}

		if err := fieldSet(field, v); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) printEnvs(s interface{}) {
	strct := structs.New(s)
	strctName := strct.Name()

	for _, field := range strct.Fields() {
		envName := strings.ToUpper(strctName) + "_" + strings.ToUpper(field.Name())
		fmt.Println("  ", envName)
	}
}

// LoadFlag creates on the fly flags based on the field names and parses them to
// load into the given pointer of struct s.
func (c *Config) LoadFlag(s interface{}) error {
	strct := structs.New(s)
	structName := strct.Name()

	f := flag.NewFlagSet(structName, flag.ExitOnError)
	// f.SetOutput(ioutil.Discard)

	for _, field := range strct.Fields() {
		name := field.Name()

		f.Var(newFieldValue(field), flagName(name), flagUsage(name))
	}

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nGenerated environment variables:\n")
		c.printEnvs(s)
		fmt.Println("")
	}

	return f.Parse(os.Args[1:])
}

// fieldSet sets field value from the given string value. It converts the
// string value in a sane way and is usefulf or environment variables or flags
// which are by nature in string types.
func fieldSet(field *structs.Field, v string) error {
	// TODO: add support for other types
	switch field.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}

		if err := field.Set(val); err != nil {
			return err
		}
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if err := field.Set(i); err != nil {
			return err
		}
	case reflect.String:
		field.Set(v)
	case reflect.Slice:
		if _, ok := field.Value().([]string); !ok {
			errors.New("can't set on non string slices")
		}

		if err := field.Set(strings.Split(v, ",")); err != nil {
			return err
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		if err := field.Set(f); err != nil {
			return err
		}
	default:
		return fmt.Errorf("multiconfig: not supported type: %s", field.Kind())
	}

	return nil
}

// fieldValue satisfies the flag.Value and flag.Getter interfaces
type fieldValue structs.Field

func newFieldValue(f *structs.Field) *fieldValue {
	fl := fieldValue(*f)
	return &fl
}

func (f *fieldValue) Set(val string) error {
	field := (*structs.Field)(f)
	return fieldSet(field, val)
}

func (f *fieldValue) String() string {
	fl := (*structs.Field)(f)
	return fmt.Sprintf("%v", fl.Value())
}

func (f *fieldValue) Get() interface{} {
	fl := (*structs.Field)(f)
	return fl.Value()
}

func flagUsage(name string) string { return fmt.Sprintf("Change value of %s.", name) }

func flagName(name string) string { return strings.ToLower(name) }
