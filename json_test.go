package multiconfig

import "testing"

func TestJSON(t *testing.T) {
	m := NewWithPath(testJSON)

	s := &Server{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testStruct(t, s, getDefaultServer())
}
