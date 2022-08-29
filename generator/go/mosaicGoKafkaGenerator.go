package gogen

import (
	"bytes"
	"embed"
	"github.com/c0olix/asyncApiCodeGen/generator"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"strings"
	"text/template"
)

var typeConversionGoMap = map[string]string{
	"int32":     "int",
	"int64":     "int64",
	"string":    "string",
	"email":     "string",
	"boolean":   "bool",
	"float":     "float32",
	"double":    "float64",
	"binary":    "[]byte",
	"date":      "time.Time",
	"date-time": "time.Time",
	"password":  "password",
}

//go:embed templates
var templateFiles embed.FS

type MosaicKafkaGoCodeGenerator struct {
	template *template.Template
	data     map[string]interface{}
	log      logrus.Logger
}

func (c MosaicKafkaGoCodeGenerator) getImports(data map[string]interface{}) []string {
	var out []string
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["publish"] != nil {
			operationProperties := operation["publish"].(map[string]interface{})
			imp := c.extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		} else if operation["subscribe"] != nil {
			operationProperties := operation["subscribe"].(map[string]interface{})
			imp := c.extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		}

	}
	return out
}

func (c MosaicKafkaGoCodeGenerator) extractImport(operationProperties map[string]interface{}) string {
	message := operationProperties["message"].(map[string]interface{})
	payload := message["payload"].(map[string]interface{})
	properties := payload["properties"].(map[string]interface{})
	for propertyName, property := range properties {
		prop := property.(map[string]interface{})
		c.log.Debugf("%v", propertyName)
		if prop["format"] != nil {
			format := prop["format"].(string)
			switch format {
			case "date-time":
				return "time"
			}
		}
	}

	return ""
}

func (c MosaicKafkaGoCodeGenerator) convertToGoType(property map[string]interface{}) string {
	switch property["type"] {
	case "object":
		if property["additionalProperties"] != nil {
			additionalProperties, ok := property["additionalProperties"].(map[string]interface{})
			if ok {
				if additionalProperties["type"] == "string" {
					return "map[string]string"
				}
			}
			_, ok = property["additionalProperties"].(bool)
			if ok {
				return property["title"].(string)
			}
		} else {
			return property["title"].(string)
		}
	case "array":
		typ := "[]"
		if property["items"] != nil {
			items := property["items"].(map[string]interface{})
			if items["format"] != nil {
				typ = typ + typeConversionGoMap[items["format"].(string)]
			} else if items["type"] == "object" {
				return items["title"].(string)
			} else {
				typ = typ + typeConversionGoMap[items["type"].(string)]
			}
			return typ
		}
	default:
		if property["format"] != nil {
			return typeConversionGoMap[property["format"].(string)]
		} else {
			return typeConversionGoMap[property["type"].(string)]
		}
	}
	return ""
}

func NewMosaicKafkaGoCodeGenerator(asyncApiSpecPath string) (*MosaicKafkaGoCodeGenerator, error) {
	goKafkaGenerator := MosaicKafkaGoCodeGenerator{}
	fns := template.FuncMap{
		"getImports":      goKafkaGenerator.getImports,
		"getMessages":     generator.GetMessages,
		"getObjects":      generator.GetNestedObjects,
		"getItemObjects":  generator.GetItemObjects,
		"convertToGoType": goKafkaGenerator.convertToGoType,
		"lower":           strings.ToLower,
		"camel":           strcase.ToCamel,
		"checkRequired":   generator.CheckRequired,
	}
	spec, err := generator.LoadAsyncApiSpecWithParser(asyncApiSpecPath)
	if err != nil {
		return nil, err
	}
	goKafkaGenerator.data = spec

	tmpl := template.Must(template.New("mosaic-kafka-go-code.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-go-code.tmpl"))
	goKafkaGenerator.template = tmpl
	return &goKafkaGenerator, nil
}

func (c MosaicKafkaGoCodeGenerator) Generate() ([]byte, error) {
	var tpl bytes.Buffer
	err := c.template.Execute(&tpl, c.data)
	if err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}