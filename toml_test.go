package multiconfig

import "testing"

func TestToml(t *testing.T) {
	m := NewWithPath(testTOML)

	s := &Server{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}
