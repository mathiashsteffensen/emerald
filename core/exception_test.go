package core_test

import (
	"emerald/object"
	"fmt"
	"strings"
	"testing"
)

func TestException_kind_of(t *testing.T) {
	tests := []coreTestCase{}

	for className := range object.Classes {
		if strings.Contains(className, "Error") {
			tests = append(tests, coreTestCase{
				name:     fmt.Sprintf("%s#kind_of?(Exception) == true", className),
				input:    fmt.Sprintf("%s.new.kind_of?(Exception)", className),
				expected: true,
			})
		}
	}

	runCoreTests(t, tests)
}
