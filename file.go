package multiconfig

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	// ErrPathNotSet states that given path to file loader is empty
	ErrPathNotSet = errors.New("config path is not set")

	// ErrFileNotFound states that given file is not exists
	ErrFileNotFound = errors.New("config file not found")
)

// TOMLLoader satisifies the loader interface. It loads the configuration from
// the given toml file or Reader
type TOMLLoader struct {
	Path   string
	Reader io.Reader
}

// Load loads the source into the config defined by struct s
// Defaults to using the Reader if provided, otherwise tries to read from the
// file
func (t *TOMLLoader) Load(s interface{}) error {
	var r io.Reader
	if t.Reader != nil {
		r = t.Reader
	} else {
		file, err := getConfig(t.Path)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	}
	if _, err := toml.DecodeReader(r, s); err != nil {
		return err
	}

	return nil
}

// JSONLoader satisifies the loader interface. It loads the configuration from
// the given json file or Reader
type JSONLoader struct {
	Path   string
	Reader io.Reader
}

// Load loads the source into the config defined by struct s
// Defaults to using the Reader if provided, otherwise tries to read from the
// file
func (j *JSONLoader) Load(s interface{}) error {
	var r io.Reader
	if j.Reader != nil {
		r = j.Reader
	} else {
		file, err := getConfig(j.Path)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	}

	return json.NewDecoder(r).Decode(s)
}

func getConfig(path string) (*os.File, error) {
	if path == "" {
		return nil, ErrPathNotSet
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(pwd, path)

	// check if file with combined path is exists(relative path)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		return os.Open(configPath)
	}

	// check if file is exists it self
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return os.Open(path)
	}

	return nil, ErrFileNotFound
}
