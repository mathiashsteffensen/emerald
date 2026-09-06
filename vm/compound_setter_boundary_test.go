package vm

import (
	"fmt"
	"strings"
	"testing"
)

func TestCompoundSetterFrameBoundary(t *testing.T) {
	for _, keywords := range []bool{false, true} {
		for _, count := range []int{252, 253, 254} {
			if keywords && count == 254 {
				continue // The setter's positional and keyword arguments must total at most 255.
			}
			params, args := make([]string, count), make([]string, count)
			for i := range params {
				params[i], args[i] = fmt.Sprintf("p%d", i), fmt.Sprint(i)
			}
			getterParams := strings.Join(params, ",")
			setterParams := getterParams + ",rhs"
			callArgs := strings.Join(args, ",")
			keywordValue := "nil"
			if keywords {
				getterParams += ",key:"
				setterParams += ",key:"
				callArgs += ",key: 99"
				keywordValue = "key"
			}
			setup := fmt.Sprintf(`class Box
				def value(%s); $gets += 1; $current; end
				def value=(%s); $sets += 1; $received = [p0, p%d, rhs, %s]; -1; end
			end
			box = Box.new; $gets = 0; $sets = 0; $received = nil`, getterParams, setterParams, count-1, keywordValue)
			for _, op := range []string{"||", "&&"} {
				for _, write := range []bool{false, true} {
					current := "42"
					if (op == "||" && write) || (op == "&&" && !write) {
						current = "false"
					}
					var result any = 42
					if current == "false" {
						result = false
					}
					var received any
					sets := 0
					if write {
						result, sets = 7, 1
						var key any
						if keywords {
							key = 99
						}
						received = []any{0, count - 1, 7, key}
					}
					runVmTests(t, []vmTestCase{{
						name:     fmt.Sprintf("args=%d/keywords=%t/%s/write=%t", count, keywords, op, write),
						input:    fmt.Sprintf("$current = %s; result = (box.value(%s) %s= 7); [result, $gets, $sets, $received]", current, callArgs, op),
						expected: []any{result, 1, sets, received},
					}}, setup)
				}
			}
		}
	}
}
