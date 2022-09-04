package cmd

import (
	"fmt"
	goGen "github.com/c0olix/asyncApiCodeGen/generator/go"
	javaGen "github.com/c0olix/asyncApiCodeGen/generator/java"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
		inputFlag, err := cmd.Flags().GetString("input")
		if err != nil {
			logger.Fatalf("Unable to get input flag: %v", err)
		} else if inputFlag == "" {
			logger.Fatal("Unable to get input flag: empty input location found: \"\"")
		}

		outputFlag, err := cmd.Flags().GetString("output")
		if err != nil {
			logger.Fatalf("Unable to get output flag: %v", err)
		} else if outputFlag == "" {
			logger.Fatal("Unable to get output flag: empty output location found: \"\"")
		}

		createDirFlag, err := cmd.Flags().GetBool("createDir")
		if err != nil {
			logger.Fatalf("Unable to get createDir flag: %v", err)
		}

		langFlag, err := cmd.Flags().GetString("lang")
		if err != nil {
			logger.Fatalf("Unable to get language flag: %v", err)
		}
		flavorFlag, err := cmd.Flags().GetString("flavor")
		if err != nil {
			logger.Fatalf("Unable to get flavor flag: %v", err)
		}

		packageFlag, err := cmd.Flags().GetString("packageName")
		if err != nil {
			logger.Fatalf("Unable to get package flag: %v", err)
		} else if packageFlag == "" {
			logger.Fatal("Unable to get packageFlag flag: empty package name found: \"\"")
		}

		if langFlag == "go" {
			generateMosaicKafkaGoCode(inputFlag, outputFlag, createDirFlag, packageFlag, flavorFlag)
		} else if langFlag == "java" {
			generateMosaicKafkaJavaCode(inputFlag, outputFlag, createDirFlag, packageFlag, flavorFlag)
		} else {
			logger.Fatalf("unsupported flags given")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("lang", "l", "go", "What kind of code should be generated?")
	generateCmd.Flags().StringP("flavor", "f", "", "Which (if) flavor should be used?")
	generateCmd.Flags().BoolP("createDir", "c", false, "Should directory be created if not present (recursive)?")
	generateCmd.Flags().StringP("output", "o", "", "Where should the generated code saved to? Attention: Go=File, Java=Dir!")
	generateCmd.Flags().StringP("packageName", "p", "", "Which package name should the generated code have?")
}

func generateMosaicKafkaGoCode(path string, out string, createDir bool, packageName string, flavor string) {
	if createDir {
		file := filepath.Dir(out)
		err := os.MkdirAll(file, os.ModePerm)
		if err != nil {
			logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to create output folder: %v", err)
		}
	}
	generator, err := goGen.NewMosaicKafkaGoCodeGenerator(path, packageName, logger)
	if err != nil {
		logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to generate code: %v", err)
	}
	output, err := generator.Generate(flavor)
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

func generateMosaicKafkaJavaCode(path string, out string, createDir bool, packageName string, flavor string) {
	if createDir {
		err := os.MkdirAll(out, os.ModePerm)
		if err != nil {
			logger.WithField("stack", fmt.Sprintf("%+v", err)).Fatalf("unable to create output folder: %v", err)
		}
	}
	gen, err := javaGen.NewMosaicKafkaJavaCodeGenerator(path, packageName, logger)
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
