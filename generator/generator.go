package generator

import (
	"embed"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Generator interface {
	Generate(asyncApiSpecPath string, out string) (string, error)
}

//go:embed templates
var templateFiles embed.FS

func loadAsyncApiSpec(asyncApiSpecPath string) asyncApiSpec {
	yamlFile, err := os.ReadFile(asyncApiSpecPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	spec := asyncApiSpec{}
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return spec
}
