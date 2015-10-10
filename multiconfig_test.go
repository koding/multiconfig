package multiconfig

import (
	"log"
	"testing"
	"time"
)

type (
	Server struct {
		Name       string `required:"true"`
		Port       int    `default:"6060"`
		ID         int64
		Labels     []int
		Enabled    bool
		Users      []string
		Postgres   Postgres
		unexported string
		Interval   time.Duration

		F00 int               `default:"1"`
		F01 int8              `default:"2"`
		F02 int16             `default:"3"`
		F03 int32             `default:"4"`
		F04 int64             `default:"5"`
		F05 uint32            `default:"6"`
		F06 float32           `default:"7.3"`
		F07 *int16            `default:"8"`
		F08 *float64          `default:"9"`
		F09 *string           `default:"ankara"`
		F10 *[]string         `default:"tr,en"`
		F11 *time.Duration    `default:"6s"`
		F12 *[]*time.Duration `default:"1ms,2m,3h"`
		F13 Int               `default:"12"`
		F14 time.Time         `default:"2015-11-05T08:15:30-05:00"`
		F15 CustomStruct
		F16 *CustomStruct
	}
	Int          int
	CustomStruct struct {
		F00 string    `default:"turkey"`
		F01 []*uint16 `default:"10,20"`
	}

	// Postgres holds Postgresql database related configuration
	Postgres struct {
		Enabled           bool
		Port              int      `required:"true" customRequired:"yes"`
		Hosts             []string `required:"true"`
		DBName            string   `default:"configdb"`
		AvailabilityRatio float64
		unexported        string
	}
)

type FlattenedServer struct {
	Postgres Postgres
}

type CamelCaseServer struct {
	AccessKey         string
	Normal            string
	DBName            string `default:"configdb"`
	AvailabilityRatio float64
}

var (
	testTOML = "testdata/config.toml"
	testJSON = "testdata/config.json"
)

func getDefaultServer() *Server {
	f07, f08 := int16(8), float64(9)
	f09, f10 := "ankara", []string{"tr", "en"}
	f12 := 6 * time.Second
	f12_0, f12_1, f12_2 := time.Millisecond, 2*time.Minute, 3*time.Hour
	f15_f01_0, f15_f01_1 := uint16(10), uint16(20)

	date := time.Time{}
	if err := date.UnmarshalText([]byte("2015-11-05T08:15:30-05:00")); err != nil {
		log.Fatalln(err)
	}

	return &Server{
		Name:     "koding",
		Port:     6060,
		Enabled:  true,
		ID:       1234567890,
		Labels:   []int{123, 456},
		Users:    []string{"ankara", "istanbul"},
		Interval: 10 * time.Second,
		Postgres: Postgres{
			Enabled:           true,
			Port:              5432,
			Hosts:             []string{"192.168.2.1", "192.168.2.2", "192.168.2.3"},
			DBName:            "configdb",
			AvailabilityRatio: 8.23,
		},

		F00: 1,
		F01: int8(2),
		F02: int16(3),
		F03: int32(4),
		F04: int64(5),
		F05: uint32(6),
		F06: float32(7.3),
		F07: &f07,
		F08: &f08,
		F09: &f09,
		F10: &f10,
		F11: &f12,
		F12: &[]*time.Duration{&f12_0, &f12_1, &f12_2},
		F13: 12,
		F14: date,
		F15: CustomStruct{
			F00: "turkey",
			F01: []*uint16{&f15_f01_0, &f15_f01_1},
		},
		F16: &CustomStruct{
			F00: "turkey",
			F01: []*uint16{&f15_f01_0, &f15_f01_1},
		},
	}
}

func getDefaultCamelCaseServer() *CamelCaseServer {
	return &CamelCaseServer{
		AccessKey:         "123456",
		Normal:            "normal",
		DBName:            "configdb",
		AvailabilityRatio: 8.23,
	}
}

func TestNewWithPath(t *testing.T) {
	var _ Loader = NewWithPath(testTOML)
}

func TestLoad(t *testing.T) {
	m := NewWithPath(testTOML)

	s := new(Server)
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}

func TestDefaultLoader(t *testing.T) {
	m := New()

	s := new(Server)
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	if err := m.Validate(s); err != nil {
		t.Error(err)
	}
	testStruct(t, s, getDefaultServer())

	s.Name = ""
	if err := m.Validate(s); err == nil {
		t.Error("Name should be required")
	}
}

func testStruct(t *testing.T, s *Server, d *Server) {
	if s.Name != d.Name {
		t.Errorf("Name value is wrong: %s, want: %s", s.Name, d.Name)
	}

	if s.Port != d.Port {
		t.Errorf("Port value is wrong: %d, want: %d", s.Port, d.Port)
	}

	if s.Enabled != d.Enabled {
		t.Errorf("Enabled value is wrong: %t, want: %t", s.Enabled, d.Enabled)
	}

	if s.Interval != d.Interval {
		t.Errorf("Interval value is wrong: %v, want: %v", s.Interval, d.Interval)
	}

	if s.ID != d.ID {
		t.Errorf("ID value is wrong: %v, want: %v", s.ID, d.ID)
	}

	if len(s.Labels) != len(d.Labels) {
		t.Errorf("Labels value is wrong: %d, want: %d", len(s.Labels), len(d.Labels))
	} else {
		for i, label := range d.Labels {
			if s.Labels[i] != label {
				t.Errorf("Label is wrong for index: %d, label: %s, want: %s", i, s.Labels[i], label)
			}
		}
	}

	if len(s.Users) != len(d.Users) {
		t.Errorf("Users value is wrong: %d, want: %d", len(s.Users), len(d.Users))
	} else {
		for i, user := range d.Users {
			if s.Users[i] != user {
				t.Errorf("User is wrong for index: %d, user: %s, want: %s", i, s.Users[i], user)
			}
		}
	}

	// Explicitly state that Enabled should be true, no need to check
	// `x == true` infact.
	if s.Postgres.Enabled != d.Postgres.Enabled {
		t.Errorf("Postgres enabled is wrong %t, want: %t", s.Postgres.Enabled, d.Postgres.Enabled)
	}

	if s.Postgres.Port != d.Postgres.Port {
		t.Errorf("Postgres Port value is wrong: %d, want: %d", s.Postgres.Port, d.Postgres.Port)
	}

	if s.Postgres.DBName != d.Postgres.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.Postgres.DBName, d.Postgres.DBName)
	}

	if s.Postgres.AvailabilityRatio != d.Postgres.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.Postgres.AvailabilityRatio, d.Postgres.AvailabilityRatio)
	}

	if len(s.Postgres.Hosts) != len(d.Postgres.Hosts) {
		// do not continue testing if this fails, because others is depending on this test
		t.Fatalf("Hosts len is wrong: %v, want: %v", s.Postgres.Hosts, d.Postgres.Hosts)
	}

	for i, host := range d.Postgres.Hosts {
		if s.Postgres.Hosts[i] != host {
			t.Fatalf("Hosts number %d is wrong: %v, want: %v", i, s.Postgres.Hosts[i], host)
		}
	}

	testStruct2(t, s, d)
}

func testStruct2(t *testing.T, s *Server, d *Server) {
	if d.F00 != s.F00 {
		t.Errorf("expected: %v, got: %v", d.F00, s.F00)
	}

	if d.F01 != s.F01 {
		t.Errorf("expected: %v, got: %v", d.F01, s.F01)
	}

	if d.F02 != s.F02 {
		t.Errorf("expected: %v, got: %v", d.F02, s.F02)
	}

	if d.F03 != s.F03 {
		t.Errorf("expected: %v, got: %v", d.F03, s.F03)
	}

	if d.F04 != s.F04 {
		t.Errorf("expected: %v, got: %v", d.F04, s.F04)
	}

	if d.F05 != s.F05 {
		t.Errorf("expected: %v, got: %v", d.F05, s.F05)
	}

	if d.F06 != s.F06 {
		t.Errorf("expected: %v, got: %v", d.F06, s.F06)
	}

	if *d.F07 != *s.F07 {
		t.Errorf("expected: %v, got: %v", *d.F07, *s.F07)
	}

	if *d.F08 != *s.F08 {
		t.Errorf("expected: %v, got: %v", *d.F08, *s.F08)
	}

	if *d.F09 != *s.F09 {
		t.Errorf("expected: %v, got: %v", *d.F09, *s.F09)
	}

	if len(*d.F10) != len(*s.F10) || (*d.F10)[0] != (*s.F10)[0] || (*d.F10)[1] != (*s.F10)[1] {
		t.Errorf("expected: %v, got: %v", *d.F10, *s.F10)
	}

	if *d.F11 != *s.F11 {
		t.Errorf("expected: %v, got: %v", d.F11, s.F11)
	}

	if len(*d.F12) != len(*s.F12) || (*(*d.F12)[0]) != (*(*s.F12)[0]) ||
		(*(*d.F12)[1]) != (*(*s.F12)[1]) || (*(*d.F12)[2]) != (*(*s.F12)[2]) {
		t.Errorf("expected: %v, got: %v", d.F12, s.F12)
	}

	if d.F13 != s.F13 {
		t.Errorf("expected: %v, got: %v", d.F13, s.F13)
	}

	if !d.F14.Equal(s.F14) {
		t.Errorf("expected: %v, got: %v", d.F14, s.F14)
	}

	if d.F15.F00 != s.F15.F00 {
		t.Errorf("expected: %v, got: %v", d.F15.F00, s.F15.F00)
	}

	if len(d.F15.F01) != len(s.F15.F01) {
		t.Errorf("expected: %v, got: %v", len(d.F15.F01), len(s.F15.F01))
	}

	if *d.F15.F01[0] != *s.F15.F01[0] || *d.F15.F01[1] != *s.F15.F01[1] {
		t.Errorf("expected: %v - %v, got: %v - %v", *d.F15.F01[0], *d.F15.F01[1], *s.F15.F01[0], *s.F15.F01[1])
	}

	if (*d.F16).F00 != (*s.F16).F00 {
		t.Errorf("expected: %v, got: %v", d.F16.F00, s.F16.F00)
	}

	if len((*d.F16).F01) != len((*s.F16).F01) {
		t.Errorf("expected: %v, got: %v", len((*d.F16).F01), len((*s.F16).F01))
	}

	if *(*d.F16).F01[0] != *(*s.F16).F01[0] || *(*d.F16).F01[1] != *(*s.F16).F01[1] {
		t.Errorf("expected: %v - %v, got: %v - %v", *(*d.F16).F01[0], *(*d.F16).F01[1], *(*s.F16).F01[0], *(*s.F16).F01[1])
	}
}

func testFlattenedStruct(t *testing.T, s *FlattenedServer, d *Server) {
	// Explicitly state that Enabled should be true, no need to check
	// `x == true` infact.
	if s.Postgres.Enabled != d.Postgres.Enabled {
		t.Errorf("Postgres enabled is wrong %t, want: %t", s.Postgres.Enabled, d.Postgres.Enabled)
	}

	if s.Postgres.Port != d.Postgres.Port {
		t.Errorf("Postgres Port value is wrong: %d, want: %d", s.Postgres.Port, d.Postgres.Port)
	}

	if s.Postgres.DBName != d.Postgres.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.Postgres.DBName, d.Postgres.DBName)
	}

	if s.Postgres.AvailabilityRatio != d.Postgres.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.Postgres.AvailabilityRatio, d.Postgres.AvailabilityRatio)
	}

	if len(s.Postgres.Hosts) != len(d.Postgres.Hosts) {
		// do not continue testing if this fails, because others is depending on this test
		t.Fatalf("Hosts len is wrong: %v, want: %v", s.Postgres.Hosts, d.Postgres.Hosts)
	}

	for i, host := range d.Postgres.Hosts {
		if s.Postgres.Hosts[i] != host {
			t.Fatalf("Hosts number %d is wrong: %v, want: %v", i, s.Postgres.Hosts[i], host)
		}
	}
}

func testCamelcaseStruct(t *testing.T, s *CamelCaseServer, d *CamelCaseServer) {
	if s.AccessKey != d.AccessKey {
		t.Errorf("AccessKey is wrong: %s, want: %s", s.AccessKey, d.AccessKey)
	}

	if s.Normal != d.Normal {
		t.Errorf("Normal is wrong: %s, want: %s", s.Normal, d.Normal)
	}

	if s.DBName != d.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.DBName, d.DBName)
	}

	if s.AvailabilityRatio != d.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.AvailabilityRatio, d.AvailabilityRatio)
	}

}

type InvalidSliceValueStruct struct {
	F00 []time.Duration `default:"5s,5x"`
}

func TestInvalidSliceValue(t *testing.T) {
	s := InvalidSliceValueStruct{}
	if err := New().Load(&s); err.Error() != "multiconfig: field 'F00' index '1' conversion err: time: unknown unit x in duration 5x" {
		t.Errorf("multiconfig: field 'F00' index '1' conversion err: time: unknown unit x in duration 5x, got: %s", err)
	}
	if nil != s.F00 {
		t.Errorf("expected: %v, got: %v", nil, s.F00)
	}
}
