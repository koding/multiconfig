package multiconfig

import (
	"strings"
	"testing"

	"github.com/ilgooz/structs"
)

func TestFlag(t *testing.T) {
	m := &FlagLoader{}
	s := &Server{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}

func TestFlagWithPrefix(t *testing.T) {
	const prefix = "Prefix"

	m := FlagLoader{Prefix: prefix}
	s := &Server{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, prefix)

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}

func TestFlattenFlags(t *testing.T) {
	m := FlagLoader{
		Flatten: true,
	}
	s := &FlattenedServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlattenedStruct(t, s, getDefaultServer())
}

func TestCamelcaseFlags(t *testing.T) {
	m := FlagLoader{
		CamelCase: true,
	}
	s := &CamelCaseServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testCamelcaseStruct(t, s, getDefaultCamelCaseServer())
}

func TestFlattenAndCamelCaseFlags(t *testing.T) {
	m := FlagLoader{
		Flatten:   true,
		CamelCase: true,
	}
	s := &FlattenedServer{}

	// get flags
	args := getFlags(t, "FlattenedCamelCaseServer", "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlattenedStruct(t, s, getDefaultServer())
}

// getFlags returns a slice of arguments that can be passed to flag.Parse()
func getFlags(t *testing.T, structName, prefix string) []string {
	if structName == "" {
		t.Fatal("struct name can not be empty")
	}

	var flags map[string]string
	switch structName {
	case "Server":
		flags = map[string]string{
			"-name":                       "koding",
			"-port":                       "6060",
			"-enabled":                    "",
			"-users":                      "ankara,istanbul",
			"-interval":                   "10s",
			"-id":                         "1234567890",
			"-labels":                     "123,456",
			"-postgres-enabled":           "",
			"-postgres-port":              "5432",
			"-postgres-hosts":             "192.168.2.1,192.168.2.2,192.168.2.3",
			"-postgres-dbname":            "configdb",
			"-postgres-availabilityratio": "8.23",

			"-f00":     "1",
			"-f01":     "2",
			"-f02":     "3",
			"-f03":     "4",
			"-f04":     "5",
			"-f05":     "6",
			"-f06":     "7.3",
			"-f07":     "8",
			"-f08":     "9",
			"-f09":     "ankara",
			"-f10":     "tr,en",
			"-f11":     "6s",
			"-f12":     "1ms,2m,3h",
			"-f13":     "12",
			"-f14":     "2015-11-05T08:15:30-05:00",
			"-f15-f00": "turkey",
			"-f15-f01": "10,20",
			"-f16-f00": "turkey",
			"-f16-f01": "10,20",
		}
	case "FlattenedServer":
		flags = map[string]string{
			"--enabled":           "",
			"--port":              "5432",
			"--hosts":             "192.168.2.1,192.168.2.2,192.168.2.3",
			"--dbname":            "configdb",
			"--availabilityratio": "8.23",
		}
	case "FlattenedCamelCaseServer":
		flags = map[string]string{
			"--enabled":            "",
			"--port":               "5432",
			"--hosts":              "192.168.2.1,192.168.2.2,192.168.2.3",
			"--db-name":            "configdb",
			"--availability-ratio": "8.23",
		}
	case "CamelCaseServer":
		flags = map[string]string{
			"--access-key":         "123456",
			"--normal":             "normal",
			"--db-name":            "configdb",
			"--availability-ratio": "8.23",
		}
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
