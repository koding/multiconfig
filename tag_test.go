package multiconfig

import "testing"

func TestDefaultValues(t *testing.T) {
	m := &TagLoader{}
	s := new(Server)
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	d := getDefaultServer()

	if s.Port != d.Port {
		t.Errorf("Port value is wrong: %d, want: %d", s.Port, d.Port)
	}

	if s.Postgres.DBName != d.Postgres.DBName {
		t.Errorf("Postgres DBName value is wrong: %s, want: %s", s.Postgres.DBName, d.Postgres.DBName)
	}

	testStruct2(t, s, d)
}
