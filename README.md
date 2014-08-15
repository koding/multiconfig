# Multiconfig [![GoDoc](https://godoc.org/github.com/koding/multiconfig?status.svg)](http://godoc.org/github.com/koding/multiconfig) [![Build Status](https://travis-ci.org/koding/multiconfig.svg)](https://travis-ci.org/koding/multiconfig) 

Load configuration from multiple sources. Multiconfig makes loading/parsing
from different configuration sources an easy task. The problem with any app is
that with time there are many options how to populate a set of configs.
Multiconfig makes it easy by dynamically creating all necessary options.
Checkout the example below to see it in action.

## Install

```bash
go get github.com/koding/multiconfig
```

## Usage and Examples

Lets define and declare a struct

```go
type Server struct {
	Name        string 
	Port          int
	Enabled     bool
	Rsers       []string 
}
```

Load the configuration :

```go
serverConf := new(Server)

// Create a new constructor without or with an initial config file
m := multiconfig.New()
m := multiconfig.NewWithPath("config.toml") // supports TOML and JSON


// Now populated the serverConf struct
err := m.Load(serverConf)

// Panic's if config cannot be loaded.
m.MustLoad(serverConf) 

```

Now run your app:

```sh
# Loads from config.toml 
$ app 

# Override any config easily with environment variables, environment variables
# are automatically generated in the form of STRUCTNAME_FIELDNAME
$ SERVER_PORT=4000 app 


# Or pass via flag. Flags are also automatically generated based on the field
# name
$ app  -port 4000
```

## TODO

* Implement --help that automatically prints the environment variables and flags

## License

The MIT License (MIT) - see LICENSE.md for more details
