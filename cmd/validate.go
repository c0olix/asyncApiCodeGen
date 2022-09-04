package cmd

import (
	"github.com/c0olix/asyncApiCodeGen/generator"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate given asyncApiSpec",
	Long:  `Validate given asyncApiSpec.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := generator.LoadAsyncApiSpecWithParser(args[0])
		if err != nil {
			logger.Fatalf("AsyncApi is invalid: %v", err)
		}
		logger.Info("AsyncApi is valid")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
