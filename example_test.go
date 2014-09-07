package multiconfig

import (
	"fmt"
	"os"
)

func ExampleEnvironmentLoader() {
	// Our struct which is used for configuration
	type ServerConfig struct {
		Name     string
		Port     int
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	// Assume those values defined before running the Loader
	os.Setenv("SERVERCONFIG_NAME", "koding")
	os.Setenv("SERVERCONFIG_PORT", "6060")

	// Instantiate loader
	l := &EnvironmentLoader{}

	s := &ServerConfig{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Here is our little config
	// Host--> koding
	// Port--> 6060
}

func ExampleTOMLLoader() {
	// Instantiate loader
	l := &TOMLLoader{Path: testTOML}

	s := &Server{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Here is our little config
	// Host--> koding
	// Port--> 6060
}

func ExampleJSONLoader() {
	// Instantiate loader
	l := &JSONLoader{Path: testJSON}

	s := &Server{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Here is our little config
	// Host--> koding
	// Port--> 6060
}

func ExampleMultiLoader() {
	os.Setenv("S_HOST", "koding")
	os.Setenv("S_PORT", "6060")

	// Instantiate loaders
	f := &FlagLoader{}
	e := &EnvironmentLoader{}

	l := MultiLoader(f, e)

	s := &Server{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Here is our little config
	// Host--> koding
	// Port--> 6060
}
