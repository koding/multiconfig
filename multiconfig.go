package multiconfig

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ilgooz/structs"
)

// Loader loads the configuration from a source. The implementer of Loader is
// responsible of setting the default values of the struct.
type Loader interface {
	// Load loads the source into the config defined by struct s
	Load(s interface{}) error
}

// DefaultLoader implements the Loader interface. It initializes the given
// pointer of struct s with configuration from the default sources. The order
// of load is TagLoader, FileLoader, EnvLoader and lastly FlagLoader. An error
// in any step stops the loading process. Each step overrides the previous
// step's config (i.e: defining a flag will override previous environment or
// file config). To customize the order use the individual load functions.
type DefaultLoader struct {
	Loader
	Validator
}

// NewWithPath returns a new instance of Loader to read from the given
// configuration file.
func NewWithPath(path string) *DefaultLoader {
	loaders := []Loader{}

	// Read default values defined via tag fields "default"
	loaders = append(loaders, &TagLoader{})

	// Choose what while is passed
	if strings.HasSuffix(path, "toml") {
		loaders = append(loaders, &TOMLLoader{Path: path})
	}

	if strings.HasSuffix(path, "json") {
		loaders = append(loaders, &JSONLoader{Path: path})
	}

	e := &EnvironmentLoader{}
	f := &FlagLoader{}

	loaders = append(loaders, e, f)
	loader := MultiLoader(loaders...)

	d := &DefaultLoader{}
	d.Loader = loader
	d.Validator = MultiValidator(&RequiredValidator{})
	return d
}

// New returns a new instance of DefaultLoader without any file loaders.
func New() *DefaultLoader {
	loader := MultiLoader(
		&TagLoader{},
		&EnvironmentLoader{},
		&FlagLoader{},
	)

	d := &DefaultLoader{}
	d.Loader = loader
	d.Validator = MultiValidator(&RequiredValidator{})
	return d
}

// MustLoadWithPath loads with the DefaultLoader settings and from the given
// Path. It exits if the config cannot be parsed.
func MustLoadWithPath(path string, conf interface{}) {
	d := NewWithPath(path)
	d.MustLoad(conf)
}

// MustLoad loads with the DefaultLoader settings. It exits if the config
// cannot be parsed.
func MustLoad(conf interface{}) {
	d := New()
	d.MustLoad(conf)
}

// MustLoad is like Load but panics if the config cannot be parsed.
func (d *DefaultLoader) MustLoad(conf interface{}) {
	if err := d.Load(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// we at koding, believe having sane defaults in our system, this is the
	// reason why we have default validators in DefaultLoader. But do not cause
	// nil pointer panics if one uses DefaultLoader directly.
	if d.Validator != nil {
		d.MustValidate(conf)
	}
}

// MustValidate validates the struct. It exits with status 1 if it can't
// validate.
func (d *DefaultLoader) MustValidate(conf interface{}) {
	if err := d.Validate(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

// fieldSet sets field value from the given string value. It converts the
// string value in a sane way and is usefulf or environment variables or flags
// which are by nature in string types.
func fieldSet(field *structs.Field, s string) error {
	if err := field.Settable(); err != nil {
		return err
	}

	return setValue(field.ReflectValue(), s, field.Name())
}

func setValue(v reflect.Value, s, name string) error {
	t := v.Type()
	k := v.Kind()

	switch k {
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()

		return setValue(v, s, name)
	case reflect.Struct:
		if v.CanAddr() {
			v = v.Addr()
		}

		if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(s))
		}

		return fmt.Errorf("multiconfig: converter not found for %v", t)
	case reflect.Slice:
		ss := strings.Split(s, ",")

		if len(ss) == 0 {
			return nil
		}

		sc := reflect.MakeSlice(t, 0, 0)

		for i, s := range ss {
			e := reflect.Indirect(reflect.New(t.Elem()))
			if err := setValue(e, s, name); err != nil {
				return fmt.Errorf("multiconfig: field '%s' index '%d' conversion err: %s", name, i, err)
			}
			sc = reflect.Append(sc, e)
		}

		v.Set(sc)
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		val, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		v.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64

		switch v.Interface().(type) {
		case time.Duration:
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			n = int64(d)
		default:
			var err error
			n, err = strconv.ParseInt(s, 10, t.Bits())
			if err != nil {
				return err
			}
		}

		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return err
		}

		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return err
		}

		v.SetFloat(n)
	default:
		return fmt.Errorf("multiconfig: field '%s' has unsupported type: %s", name, k)
	}

	return nil
}

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func implementsTextUnmarshaler(field *structs.Field) bool {
	t := field.Type()
	if t.Kind() == reflect.Ptr {
		return t.Implements(textUnmarshalerType)
	}
	return reflect.PtrTo(t).Implements(textUnmarshalerType)
}
