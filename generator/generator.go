package generator

import (
	"bytes"
	"go/format"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"text/template"
)

type Generator interface {
	Generate(asyncApiSpecPath string, out string) (string, error)
}

type ServerCodeGenerator struct {
	template *template.Template
}

func NewServerCodeGenerator() ServerCodeGenerator {
	tmpl := template.Must(template.ParseFiles("generator/templates/server-code.tmpl"))
	return ServerCodeGenerator{
		template: tmpl,
	}
}

func (c ServerCodeGenerator) Generate(asyncApiSpecPath string, out string) (string, error) {
	spec := c.loadAsyncApiSpec(asyncApiSpecPath)
	spec.convertToUsableStruct()
	var tpl bytes.Buffer
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	err = c.template.Execute(&tpl, spec)
	if err != nil {
		return "", err
	}
	p, err := format.Source(tpl.Bytes())
	if err != nil {
		return "", err
	}
	_, err = f.Write(p)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (c ServerCodeGenerator) loadAsyncApiSpec(asyncApiSpecPath string) asyncApiSpec {
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
