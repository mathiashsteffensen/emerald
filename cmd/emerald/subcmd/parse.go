package subcmd

import (
	"emerald/cmd/helpers"
	"emerald/log"
	"emerald/parser"
	"emerald/parser/lexer"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var ParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Run the emerald parser, outputs the stringified AST",
	Long:  "Primarily just used for debugging and benchmarking",
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			absFile, err := filepath.Abs(file)
			helpers.CheckError("Failed to make path absolute?", err)

			bytes, err := os.ReadFile(absFile)
			helpers.CheckError("error reading file", err)

			start := time.Now()

			l := lexer.New(lexer.NewInput(absFile, string(bytes)))
			p := parser.New(l)
			program := p.ParseAST()

			done := time.Since(start)

			if len(p.Errors()) != 0 {
				log.FatalF("parser error: %s\n", p.Errors()[0])
			}

			log.Debug("\n" + program.String())

			log.DebugF("Parsed file %s in %s", file, done)
		}
	},
}
