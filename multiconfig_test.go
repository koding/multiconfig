package multiconfig

import "testing"

type (
	Server struct {
		Name     string
		Port     int
		Enabled  bool
		Users    []string
		Postgres Postgres
	}

	// Postgres holds Postgresql database related configuration
	Postgres struct {
		Enabled           bool
		Port              int
		Hosts             []string
		DBName            string
		AvailabilityRatio float64
	}
)

var (
	testTOML = "testdata/config.toml"
	testJSON = "testdata/config.json"
)

func TestNewWithPath(t *testing.T) {
	var _ Loader = NewWithPath(testTOML)
}

func TestLoad(t *testing.T) {
	m := NewWithPath(testTOML)

	s := new(Server)
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	if s.Name != "Koding" {
		t.Errorf("Name value is wrong: %s, want: %s", s.Name, "Koding")
	}

	if s.Port != 6060 {
		t.Errorf("Port value is wrong: %s, want: %s", s.Port, 6060)
	}

	if !s.Enabled {
		t.Errorf("Enabled value is wrong: %s, want: %s", s.Enabled, true)
	}

	if len(s.Users) != 2 {
		t.Errorf("Users value is wrong: %s, want: %s", len(s.Users), 2)
	}
}

func TestTomlEmbeddedStruct(t *testing.T) {
	m := NewWithPath(testTOML)

	s := &Server{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEmbededStruct(t, s)
}

func TestJSONEmbeddedStruct(t *testing.T) {
	m := NewWithPath(testJSON)

	s := &Server{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEmbededStruct(t, s)
}

func testEmbededStruct(t *testing.T, s *Server) {
	// Explicitly state that Enabled should be true, no need to check
	// `x == true` infact.
	if s.Postgres.Enabled != true {
		t.Error("Enabled should be true")
	}

	port := 5432
	if s.Postgres.Port != port {
		t.Errorf("Port value is wrong: %s, want: %s", s.Postgres.Port, port)
	}

	dbName := "configdb"
	if s.Postgres.DBName != dbName {
		t.Errorf("DBName is wrong: %s, want: %s", s.Postgres.DBName, dbName)
	}

	var availabilityRatio float64 = 8.23
	if s.Postgres.AvailabilityRatio != availabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %d, want: %d", s.Postgres.AvailabilityRatio, availabilityRatio)
	}

	hosts := []string{"192.168.2.1", "192.168.2.2", "192.168.2.3"}
	if len(s.Postgres.Hosts) != 3 {
		// do not continue testing if this fails, because others is depending on this test
		t.Fatalf("Hosts len is wrong: %v, want: %v", s.Postgres.Hosts, hosts)
	}

	for i, host := range hosts {
		if s.Postgres.Hosts[i] != host {
			t.Fatalf("Hosts number %d is wrong: %v, want: %v", i, s.Postgres.Hosts[i], host)
		}
	}
}
