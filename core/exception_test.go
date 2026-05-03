package core_test

import (
	"emerald/core"
	"fmt"
	"strings"
	"testing"
)

func TestException_kind_of(t *testing.T) {
	tests := []coreTestCase{}

	rt := core.NewRuntime()
	rt.Init()

	for className := range rt.Object.NamespaceDefinitions() {
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
