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
		langFlag, err := cmd.Flags().GetString("lang")
		if err != nil {
			logger.Errorf("Unable to get language flag: %v", err)
		}
		flavorFlag, err := cmd.Flags().GetString("flavor")
		if err != nil {
			logger.Errorf("Unable to get flavor flag: %v", err)
		}
		if langFlag == "go" && flavorFlag == "mosaic" {
			generateMosaicKafkaGoCode(args[0], args[1])
		} else if langFlag == "java" && flavorFlag == "mosaic" {
			generateMosaicKafkaJavaCode(args[0], args[1])
		} else {
			logger.Fatalf("unsupported flags given")
		}

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
	generateCmd.Flags().StringP("lang", "l", "go", "What kind of code should be generated?")
	generateCmd.Flags().StringP("flavor", "f", "mosaic", "Which flavor should be used??")
}

func generateMosaicKafkaGoCode(path string, out string) {
	logger.Debug("generate called")
	gen := generator.NewMosaicKafkaGoCodeGenerator()
	out, err := gen.Generate(path, out)
	if err != nil {
		logger.Fatalf("unable to generate code: %v", err)
	}
	logger.Debug(out)
}

func generateMosaicKafkaJavaCode(path string, out string) {
	logger.Debug("generate called")
	gen := generator.NewMosaicKafkaJavaCodeGenerator()
	out, err := gen.Generate(path, out)
	if err != nil {
		logger.Fatalf("unable to generate code: %v", err)
	}
	logger.Debug(out)
}
