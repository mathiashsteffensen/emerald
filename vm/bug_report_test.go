package vm

import (
	"emerald/core"
	"emerald/object"
	"testing"
)

func TestBugReportTradingStrategyBuiltinInIterator(t *testing.T) {
	symbols := []string{"XOP", "C", "AMD", "GOOGL", "PFE", "SIRI"}
	tests := []vmTestCase{
		{
			name: "builtin method call before hash assignment in iterator",
			input: `
			  result = {}
				$symbols.each do |symbol|
				  puts("EMA:", ema(symbol, 10))
				  result[symbol] = BUY
				end
				result.to_s
			`,
			expected: "{XOP => buy, C => buy, AMD => buy, GOOGL => buy, PFE => buy, SIRI => buy}",
		},
	}

	runVmTestsWithRuntimeSetup(t, tests, func(rt *core.Runtime) {
		emeraldSymbols := make([]object.EmeraldValue, 0, len(symbols))
		for _, symbol := range symbols {
			emeraldSymbols = append(emeraldSymbols, rt.NewString(symbol))
		}
		rt.Heap.SetGlobalVariableString("$symbols", rt.NewArray(emeraldSymbols))
		rt.DefineGlobalConstant("BUY", rt.NewString("buy"))

		emaCalls := 0
		rt.DefineMethod(rt.Object, "ema", func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
			if _, err := rt.EnforceArity(args, kwargs, 2, 2); err != nil {
				return object.NewHeapObject(err)
			}

			symbolArg, err := core.EnforceArgumentType[*core.StringInstance](rt, rt.String, args[0])
			if err != nil {
				return object.NewHeapObject(err)
			}

			period, err := rt.EnforceIntegerArg(args[1])
			if err != nil {
				return object.NewHeapObject(err)
			}

			if emaCalls >= len(symbols) {
				t.Fatalf("ema called too many times")
			}
			if symbolArg.Value != symbols[emaCalls] {
				t.Fatalf("ema call %d symbol = %s, expected %s", emaCalls, symbolArg.Value, symbols[emaCalls])
			}
			if period != 10 {
				t.Fatalf("ema call %d period = %d, expected 10", emaCalls, period)
			}

			emaCalls++
			return rt.NewFloat(5.620074)
		})
	})
}

func TestBugReportUnhandledExceptionStopsBuiltinIterator(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "raised builtin method in iterator preserves original exception",
			input: `
			  result = {}
				$symbols.each do |symbol|
				  puts("EMA:", ema(symbol, 10))
				  result[symbol] = BUY
				end
				result
			`,
			expected: "error:StandardError:no stats for symbol: PATH",
		},
	}

	runVmTestsWithRuntimeSetup(t, tests, func(rt *core.Runtime) {
		rt.Heap.SetGlobalVariableString("$symbols", rt.NewArray([]object.EmeraldValue{
			rt.NewString("PBR"),
			rt.NewString("PATH"),
			rt.NewString("C"),
		}))
		rt.DefineGlobalConstant("BUY", rt.NewString("buy"))

		rt.DefineMethod(rt.Object, "ema", func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
			if _, err := rt.EnforceArity(args, kwargs, 2, 2); err != nil {
				return object.NewHeapObject(err)
			}

			symbolArg, err := core.EnforceArgumentType[*core.StringInstance](rt, rt.String, args[0])
			if err != nil {
				return object.NewHeapObject(err)
			}

			if _, err := rt.EnforceIntegerArg(args[1]); err != nil {
				return object.NewHeapObject(err)
			}

			if symbolArg.Value == "PATH" {
				return object.NewHeapObject(rt.Raise(rt.NewStandardError("no stats for symbol: PATH")))
			}

			return rt.NewFloat(2.059014)
		})
	})
}
