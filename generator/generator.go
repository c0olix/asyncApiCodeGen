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

type MosaicKafkaGoCodeGenerator struct {
	template *template.Template
}

func NewMosaicKafkaGoCodeGenerator() MosaicKafkaGoCodeGenerator {
	tmpl := template.Must(template.ParseFiles("generator/templates/mosaic-kafka-go-code.tmpl"))
	return MosaicKafkaGoCodeGenerator{
		template: tmpl,
	}
}

func (c MosaicKafkaGoCodeGenerator) Generate(asyncApiSpecPath string, out string) (string, error) {
	spec := loadAsyncApiSpec(asyncApiSpecPath)
	spec.convertToGoSpec()
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

type MosaicKafkaJavaCodeGenerator struct {
	eventClassTemplate *template.Template
}

func NewMosaicKafkaJavaCodeGenerator() MosaicKafkaJavaCodeGenerator {
	tmpl := template.Must(template.ParseFiles("generator/templates/mosaic-kafka-java-event-class.tmpl"))
	return MosaicKafkaJavaCodeGenerator{
		eventClassTemplate: tmpl,
	}
}

func (c MosaicKafkaJavaCodeGenerator) Generate(asyncApiSpecPath string, out string) (string, error) {
	spec := loadAsyncApiSpec(asyncApiSpecPath)
	spec.convertToJavaSpec()

	for _, event := range spec.Events {
		var tpl bytes.Buffer
		f, err := os.Create(out + "/" + event.Name + ".java")
		if err != nil {
			return "", err
		}
		err = c.eventClassTemplate.Execute(&tpl, event)
		if err != nil {
			return "", err
		}
		_, err = f.Write(tpl.Bytes())
		if err != nil {
			return "", err
		}
	}

	return "", nil
}
