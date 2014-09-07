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

	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Host--> koding
	// Port--> 6060
}

func ExampleTOMLLoader() {
	// Our struct which is used for configuration
	type ServerConfig struct {
		Name     string
		Port     int
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	// Instantiate loader
	l := &TOMLLoader{Path: testTOML}

	s := &ServerConfig{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Host-->", s.Name)
	fmt.Println("Users-->", s.Users)

	// Output:
	// Host--> koding
	// Users--> [ankara istanbul]
}

func ExampleJSONLoader() {
	// Our struct which is used for configuration
	type ServerConfig struct {
		Name     string
		Port     int
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	// Instantiate loader
	l := &JSONLoader{Path: testJSON}

	s := &ServerConfig{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Host-->", s.Name)
	fmt.Println("Users-->", s.Users)

	// Output:
	// Host--> koding
	// Users--> [ankara istanbul]
}

func ExampleMultiLoader() {
	// Our struct which is used for configuration
	type ServerConfig struct {
		Name     string
		Port     int
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	os.Setenv("SERVERCONFIG_HOST", "koding")
	os.Setenv("SERVERCONFIG_PORT", "6060")

	// Create a multi loader intance based on your needs
	f := &FlagLoader{}
	e := &EnvironmentLoader{}
	l := MultiLoader(f, e)

	// Load configs into our s variable from the sources above
	s := &ServerConfig{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Host--> koding
	// Port--> 6060
}
