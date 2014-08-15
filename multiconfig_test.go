package multiconfig

import "testing"

var (
	testPath = "testconfig"
)

func TestNewWithPath(t *testing.T) {
	m := NewWithPath(testPath)

	if m.Path == "" {
		t.Error("Path should be not empty")
	}
}
