/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
