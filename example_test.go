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
	os.Setenv("S_HOST", "golang.org")
	os.Setenv("S_PORT", "80")
	l := &EnvironmentLoader{}
	err := l.Load(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Host-->", s.Host)
	fmt.Println("Port-->", s.Host)
}
