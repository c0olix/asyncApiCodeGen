/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/c0olix/asyncApiCodeGen/generator"
	"github.com/c0olix/asyncApiCodeGen/logging"
	"github.com/spf13/cobra"
)

var logger = logging.NewLogger()

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "used to generate code from async api spec",
	Long: `used to generate code from async api spec. For example:
	
	asyncApiCodeGen generate bla.yaml
 `,
	Run: func(cmd *cobra.Command, args []string) {
		generateServer(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().StringP("type", "t", "server", "What kind of code should be generated? Server or Client")
}

func generateServer(path string, out string) {
	logger.Debug("generate called")
	gen := generator.NewServerCodeGenerator()
	out, err := gen.Generate(path, out)
	if err != nil {
		logger.Fatalf("unable to generate code: %v", err)
	}
	logger.Debug(out)
}
