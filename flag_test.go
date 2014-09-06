package multiconfig

import (
	"os"
	"strings"
	"testing"

	"github.com/fatih/structs"
)

func TestFlag(t *testing.T) {
	m := &FlagLoader{}
	s := &Server{}
	structName := structs.Name(s)

	// set env variables
	args := getFlags(t, structName, "")

	os.Args = args

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}

// getFlags returns a slice of arguments that can be passed to flag.Parse()
func getFlags(t *testing.T, structName, prefix string) []string {
	if structName == "" {
		t.Fatal("struct name can not be empty")
	}

	flags := map[string]string{
		"-name":                       "koding",
		"-port":                       "6060",
		"-enabled":                    "",
		"-users":                      "ankara,istanbul",
		"-postgres-enabled":           "",
		"-postgres-port":              "5432",
		"-postgres-hosts":             "192.168.2.1,192.168.2.2,192.168.2.3",
		"-postgres-dbname":            "configdb",
		"-postgres-availabilityratio": "8.23",
	}

	prefix = strings.ToLower(prefix)

	args := []string{"multiconfig-test"}
	for key, val := range flags {
		flag := key
		if prefix != "" {
			flag = "-" + prefix + key
		}

		if val == "" {
			args = append(args, flag)
		} else {
			args = append(args, flag, val)
		}
	}

	return args
}
