package core_test

import "testing"

func TestDir_pwd(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `!!(Dir.pwd =~ /emerald\/core$/)`,
			expected: true,
		},
	}

	runCoreTests(t, tests)
}

func TestDir_glob(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `Dir.glob("fixtures/*.rb")`,
			expected: []any{"fixtures/require_test.rb", "fixtures/require_test_2.rb"},
		},
		{
			input:    `Dir.glob("fixtures/*.{rb,em}")`,
			expected: []any{
				// TODO:
				// "fixtures/glob_test.em", "fixtures/require_test.rb", "fixtures/require_test_2.rb"
			},
		},
	}

	runCoreTests(t, tests)
}
