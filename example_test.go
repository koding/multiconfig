package multiconfig

import (
	"fmt"
	"os"
)

type S struct {
	Host string
	Port int
}

func ExampleEnvironmentLoader() {
	// Assume those values defined before running the Loader
	os.Setenv("S_HOST", "golang.org")
	os.Setenv("S_PORT", "80")

	// Instantiate loader
	l := &EnvironmentLoader{}
	s := &S{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}

func ExampleTOMLLoader() {
	const path = "/path/to/config.toml"

	// Instantiate loader
	l := &TOMLLoader{Path: path}

	s := &S{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}

func ExampleJSONLoader() {
	const path = "/path/to/config.json"

	// Instantiate loader
	l := &JSONLoader{Path: path}

	s := &S{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}

func ExampleFlagLoader() {
	// Instantiate loader
	l := &FlagLoader{}

	s := &S{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}

func ExampleMultiLoader() {
	// Instantiate loaders
	f := &FlagLoader{}
	e := &EnvironmentLoader{}

	l := MultiLoader(f, e)

	s := &S{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is our little config")
	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}
