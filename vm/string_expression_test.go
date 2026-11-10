package vm

import "testing"

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"", `"monkey"`, "monkey"},
		{"", `"mon" + "key"`, "monkey"},
		{"", `"mon" + "key" + "banana"`, "monkeybanana"},
		{"", `placeholder = "template"; "This is a #{placeholder}"`, "This is a template"},
		{"user defined to_s", `class Point; def to_s; "P(#{1}, #{2})"; end; end; "#{Point.new} text"`, "P(1, 2) text"},
		{"user defined to_s twice", `class Point; def to_s; "P"; end; end; point = Point.new; "a #{point} b #{point} c"`, "a P b P c"},
		{"raising to_s", `class Bad; def to_s; raise "bad"; end; end; "a #{Bad.new} b"`, "error:RuntimeError:bad"},
	}
	runVmTests(t, tests)
}
