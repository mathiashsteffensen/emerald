package emerald

import (
	"github.com/elliotchance/orderedmap/v2"
	"testing"
)

var testCases = map[string]any{
	"Hello sir": 2,
	"boop":      "bob",
}

func BenchmarkStdlibMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]any{}

		for key, value := range testCases {
			m[key] = value
		}

		for key, value := range testCases {
			if m[key] != value {
				b.Fatalf("Shitty map")
			}
		}
	}
}

func BenchmarkPackageOrderedMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := orderedmap.NewOrderedMap[string, any]()

		for key, value := range testCases {
			m.Set(key, value)
		}

		for key, value := range testCases {
			if m.GetOrDefault(key, nil) != value {
				b.Fatalf("Shitty map")
			}
		}
	}
}
