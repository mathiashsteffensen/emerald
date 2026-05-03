package main

import (
	"emerald"
	"emerald/cmd/emerald/subcmd"
	"emerald/cmd/helpers"
	"emerald/debug"
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
		debug.ExperimentalWarning()

		if logHeapUsage {
			go logHeapUsageRoutine()
		}

		defer helpers.RecoverWithStacktrace()

		engine := emerald.New()

		for _, file := range args {
			absFile, err := filepath.Abs(file)
			helpers.CheckError("Failed to make path absolute?", err)

			bytes, err := os.ReadFile(absFile)
			helpers.CheckError("error reading file", err)

			_, err = engine.EvalFile(absFile, string(bytes))

			if err != nil {
				debug.Fatal(err.Error())
			}

			debug.Shutdown()
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
		debug.FatalF("error: %s", err)
		os.Exit(1)
	}
}

func logHeapUsageRoutine() {
	m := runtime.MemStats{}

	for {
		time.Sleep(200 * time.Millisecond)

		runtime.ReadMemStats(&m)
		heapAlloc := float64(m.HeapAlloc) / 1024 / 1024 // In MB

		debug.DebugF("Heap size: %fMB", heapAlloc)
	}
}
