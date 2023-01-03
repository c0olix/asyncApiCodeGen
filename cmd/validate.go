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
		inputFlag, err := cmd.Flags().GetString("input")
		if err != nil {
			logger.Fatalf("Unable to get input flag: %v", err)
		} else if inputFlag == "" {
			logger.Fatal("Unable to get input flag: empty input location found: \"\"")
		}

		_, err = generator.LoadAsyncApiSpecWithParser(inputFlag)
		if err != nil {
			logger.Fatalf("AsyncApi is invalid: %v", err)
		}
		logger.Info("AsyncApi is valid")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
