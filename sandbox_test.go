package gronx

import (
	"strings"
	"testing"
)

func TestSegments(t *testing.T) {
	t.Run("test segment", func(t *testing.T) {
		result, err := Segments("* * * * *")
		if err != nil {
			t.Errorf(err.Error())
		}

		expect := []string{"*", "*", "*", "*", "*"}
		//if reflect.DeepEqual(result, expect) {
		if len(result) != len(expect) {
			t.Errorf("segments not expected, expected: %s, got: %s", strings.Join(expect, ", "), strings.Join(result, ", "))
		}
	})
}
