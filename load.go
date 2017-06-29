package multiconfig

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Confer interface {
	GetConf() string
}

// MustLoadInTurn load configuration from tag/json/yml/yaml/toml/env/flag(from low to heigh) and validate.
// If error, exit 2
func MustLoadInTurn(app string, conf Confer) {
	flagLoader := &FlagLoader{
		Prefix:        "",
		Flatten:       false,
		CamelCase:     true,
		EnvPrefix:     app,
		ErrorHandling: flag.ExitOnError,
	}

	for _, loader := range []Loader{
		&TagLoader{},
		&JSONLoader{Path: "conf/" + app + ".json"},
		&YAMLLoader{Path: "conf/" + app + ".yml"},
		&YAMLLoader{Path: "conf/" + app + ".yaml"},
		&TOMLLoader{Path: "conf/" + app + ".toml"},
		&EnvironmentLoader{Prefix: app, CamelCase: true},
		flagLoader,
	} {
		loader.Load(conf)
	}

	if conf.GetConf() != "" {
		path := conf.GetConf()
		loaders := []Loader{}
		if strings.HasSuffix(path, "toml") {
			loaders = append(loaders, &TOMLLoader{Path: path})
		} else if strings.HasSuffix(path, "json") {
			loaders = append(loaders, &JSONLoader{Path: path})
		} else if strings.HasSuffix(path, "yml") || strings.HasSuffix(path, "yaml") {
			loaders = append(loaders, &YAMLLoader{Path: path})
		}
		loaders = append(loaders, &FlagLoader{})
		loader := MultiLoader(loaders...)
		if err := loader.Load(conf); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	rv := RequiredValidator{}
	if err := rv.Validate(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
