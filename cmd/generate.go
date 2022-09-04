package cmd

import (
	"fmt"
	goGen "github.com/c0olix/asyncApiCodeGen/generator/go"
	javaGen "github.com/c0olix/asyncApiCodeGen/generator/java"
	"github.com/spf13/cobra"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Used to generate code from async api spec",
	Long: `Used to generate code from async api spec. First argument is the spec and the second is the path to the output. 
	
	For example:
	asyncApiCodeGen generate in_spec.yaml out.go
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
	generateCmd.Flags().StringP("lang", "l", "go", "What kind of code should be generated?")
	generateCmd.Flags().StringP("flavor", "f", "mosaic", "Which flavor should be used??")
}

func generateMosaicKafkaGoCode(path string, out string) {
	generator, err := goGen.NewMosaicKafkaGoCodeGenerator(path, logger)
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to generate code: %v", err)
	}
	output, err := generator.Generate()
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to generate code: %v", err)
	}
	f, err := os.Create(out)
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to create output file: %v", err)
	}
	_, err = f.Write(output)
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to write to output file: %v", err)
	}
}

func generateMosaicKafkaJavaCode(path string, out string) {
	err := os.MkdirAll(out, os.ModePerm)
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to create output folder: %v", err)
	}
	gen, err := javaGen.NewMosaicKafkaJavaCodeGenerator(path, logger)
	if err != nil {
		logger.Fatalf("unable to generate code: %v", err)
	}
	results, err := gen.Generate()
	if err != nil {
		logger.Fatalf("unable to generate code: %v", err)
	}
	for _, result := range results.Files {
		f, err := os.Create(out + "/" + result.Name + ".java")
		if err != nil {
			logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to create output file: %v", err)
		}
		_, err = f.Write(result.Content)
		if err != nil {
			logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to write to output file: %v", err)
		}
	}
}
