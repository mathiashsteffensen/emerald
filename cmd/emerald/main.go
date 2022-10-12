package main

import (
	"emerald/cmd/emerald/subcmd"
	"emerald/cmd/helpers"
	"emerald/compiler"
	"emerald/heap"
	"emerald/log"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/lexer"
	"emerald/vm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var profilingEnabled bool
var logHeapUsage bool

var rootCmd = &cobra.Command{
	Use:   "emerald",
	Short: "A Ruby implementation written in Go",
	Long:  "Emerald is a Ruby compiler & Virtual Machine implemented in Go",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.ExperimentalWarning()

		defer helpers.RecoverWithStacktrace()

		for _, file := range args {
			absFile, err := filepath.Abs(file)
			helpers.CheckError("Failed to make path absolute?", err)

			bytes, err := os.ReadFile(absFile)
			helpers.CheckError("error reading file", err)

			l := lexer.New(lexer.NewInput(absFile, string(bytes)))
			p := parser.New(l)
			program := p.ParseAST()

			if len(p.Errors()) != 0 {
				log.FatalF("parser error: %s\n", p.Errors()[0])
			}

			c := compiler.New()

			err = c.Compile(program)
			helpers.CheckError("Compilation failed", err)

			if logHeapUsage {
				go logHeapUsageRoutine()
			}

			machine := vm.New(absFile, c.Bytecode())
			machine.Run()

			if exception := heap.GetGlobalVariableString("$!"); exception != nil {
				exception := exception.(object.EmeraldError)
				log.FatalF("%s: %s", exception.ClassName(), exception.Message())
			}

			if machine.StackTop() != nil {
				log.InternalDebug("StackTop was not nil")
			}

			log.Shutdown()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&profilingEnabled, "profile", false, "EM_DEBUG=1 emerald --profile lib/main.rb")
	rootCmd.PersistentFlags().BoolVar(&logHeapUsage, "logHeapUsage", false, "EM_DEBUG=1 emerald --logHeapUsage lib/main.rb")

	rootCmd.AddCommand(subcmd.ParseCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.FatalF("error: %s", err)
		os.Exit(1)
	}
}

func logHeapUsageRoutine() {
	m := runtime.MemStats{}

	for {
		time.Sleep(200 * time.Millisecond)

		runtime.ReadMemStats(&m)
		heapAlloc := float64(m.HeapAlloc) / 1024 / 1024 // In MB

		log.DebugF("Heap size: %fMB", heapAlloc)
	}
}
