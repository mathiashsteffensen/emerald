package vm

import "testing"

func TestEnsure(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "runs after the body, which gives the value", input: `
			def f
				$log << :body
				1
			ensure
				$log << :ensure
				2
			end
			[f, $log]
		`, expected: []any{1, []any{":body", ":ensure"}}},
		{name: "runs on return", input: `
			def f
				return 1 if true
				$log << :not_reached
			ensure
				$log << :ensure
			end
			[f, $log]
		`, expected: []any{1, []any{":ensure"}}},
		{name: "runs after rescue, which gives the value", input: `
			def f
				raise "boom"
			rescue => error
				$log << :rescue
				error.message
			ensure
				$log << :ensure
			end
			[f, $log]
		`, expected: []any{"boom", []any{":rescue", ":ensure"}}},
		{name: "runs before an error it does not rescue reaches the caller", input: `
			def f
				raise ArgumentError, "boom"
			rescue TypeError
				$log << :wrong_rescue
			ensure
				$log << :ensure
			end
			def g
				f
			rescue ArgumentError => error
				$log << error.message
			end
			g
			$log
		`, expected: []any{":ensure", "boom"}},
		{name: "runs before an error from rescue reaches the caller", input: `
			def f
				raise "first"
			rescue
				raise "second"
			ensure
				$log << :ensure
			end
			def g
				f
			rescue => error
				$log << error.message
			end
			g
			$log
		`, expected: []any{":ensure", "second"}},
		{name: "runs when an error unwinds a native iterator", input: `
			def f
				[1].each { raise "boom" }
			ensure
				$log << :ensure
			end
			def g
				f
			rescue => error
				$log << error.message
			end
			g
			$log
		`, expected: []any{":ensure", "boom"}},
		{name: "runs in a block before its error leaves a native iterator", input: `
			def g
				[1].each do |n|
					raise "boom"
				ensure
					$log << :ensure
				end
			rescue => error
				$log << error.message
			end
			g
			$log
		`, expected: []any{":ensure", "boom"}},
		{name: "an error in ensure is not rescued by its body", input: `
			def f
				1
			rescue
				$log << :rescue
			ensure
				raise "from ensure"
			end
			def g
				f
			rescue => error
				$log << error.message
			end
			g
			$log
		`, expected: []any{"from ensure"}},
		{name: "return in ensure discards the error", input: `
			def f
				raise "lost"
			ensure
				return :ensure
			end
			f
		`, expected: ":ensure"},
		{name: "empty", input: "def f\n1\nensure\nend; f", expected: 1},
		{name: "runs after an empty body", input: "def f\nensure\n$log << :ensure\nend\n[f, $log]", expected: []any{nil, []any{":ensure"}}},
		{name: "runs on next in a block", input: `
			result = [1, 2].map do |n|
				next n * 10 if n == 1
				n
			ensure
				$log << n
			end
			[result, $log]
		`, expected: []any{[]any{10, 2}, []any{1, 2}}},
		{name: "keeps the error when nothing rescues it", input: `
			def f
				raise "boom"
			ensure
				$log << :ensure
			end
			f
		`, expected: "error:RuntimeError:boom"},
	}, "$log = []")
}
