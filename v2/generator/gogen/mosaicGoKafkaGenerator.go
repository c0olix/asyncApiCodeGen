package gogen

import (
	"bytes"
	"fmt"
	"github.com/c0olix/asyncApiCodeGen/v2/generator"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"go/format"
	"golang.org/x/exp/slices"
	"strings"
	"text/template"
)

type MosaicKafkaGoCodeGenerator struct {
	mosaicTemplate *template.Template
	normalTemplate *template.Template
	mqttTemplate   *template.Template
	data           map[string]interface{}
	log            *logrus.Logger
}

func (thiz *MosaicKafkaGoCodeGenerator) getImports(data map[string]interface{}) []string {
	var out []string
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["publish"] != nil {
			operationProperties := operation["publish"].(map[string]interface{})
			imp := thiz.extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		} else if operation["subscribe"] != nil {
			operationProperties := operation["subscribe"].(map[string]interface{})
			imp := thiz.extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		}

	}
	return out
}

func (thiz *MosaicKafkaGoCodeGenerator) extractImport(operationProperties map[string]interface{}) string {
	message := operationProperties["message"].(map[string]interface{})
	payload := message["payload"].(map[string]interface{})
	properties := payload["properties"].(map[string]interface{})
	for _, property := range properties {
		prop := property.(map[string]interface{})
		if prop["format"] != nil {
			frm := prop["format"].(string)
			switch frm {
			case "date-time":
				return "time"
			}
		}
	}

	return ""
}

func (thiz *MosaicKafkaGoCodeGenerator) convertToGoType(property map[string]interface{}) string {
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
				typ = typ + items["title"].(string)
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

func (thiz *MosaicKafkaGoCodeGenerator) validations(property map[string]interface{}, required bool) string {
	out := ` validate:"`
	if required && property["type"] != "boolean" {
		out = out + "required,"
	}
	if property["format"] != nil {
		switch property["format"].(string) {
		case "email":
			out = out + "email,"
		}
	}
	if property["minimum"] != nil {
		min := property["minimum"].(int)
		out = out + fmt.Sprintf("gte=%d,", min)
	} else if property["exclusiveMinimum"] != nil {
		min := property["exclusiveMinimum"].(int)
		out = out + fmt.Sprintf("gt=%d,", min)
	}
	if property["maximum"] != nil {
		max := property["maximum"].(int)
		out = out + fmt.Sprintf("lte=%d,", max)
	} else if property["exclusiveMaximum"] != nil {
		max := property["exclusiveMaximum"].(int)
		out = out + fmt.Sprintf("lt=%d,", max)
	}
	if property["type"] == "string" {
		if property["minLength"] != nil {
			min := property["minLength"].(int)
			out = out + fmt.Sprintf("min=%d,", min)
		}
		if property["maxLength"] != nil {
			max := property["maxLength"].(int)
			out = out + fmt.Sprintf("max=%d,", max)
		}
	}
	if property["type"] == "array" {
		if property["minItems"] != nil {
			min := property["minItems"].(int)
			out = out + fmt.Sprintf("min=%d,", min)
		}
		if property["maxItems"] != nil {
			max := property["maxItems"].(int)
			out = out + fmt.Sprintf("max=%d,", max)
		}
		if property["uniqueItems"] != nil {
			if property["uniqueItems"].(bool) {
				out = out + "unique,"
			}
		}
	}
	out = out[:len(out)-1] + `"`
	if out != ` validate:"` {
		return out
	}
	return ""
}

func NewMosaicKafkaGoCodeGenerator(asyncApiSpecPath string, packageName string, log *logrus.Logger) (*MosaicKafkaGoCodeGenerator, error) {
	goKafkaGenerator := MosaicKafkaGoCodeGenerator{
		log: log,
	}
	fns := template.FuncMap{
		"getImports":      goKafkaGenerator.getImports,
		"getMessages":     generator.GetMessages,
		"getObjects":      generator.GetNestedObjects,
		"getItemObjects":  generator.GetItemObjects,
		"convertToGoType": goKafkaGenerator.convertToGoType,
		"lower":           strings.ToLower,
		"camel":           strcase.ToCamel,
		"lowerCamel":      strcase.ToLowerCamel,
		"checkRequired":   generator.CheckRequired,
		"validations":     goKafkaGenerator.validations,
		"hasProducer":     generator.HasProducer,
		"hasConsumer":     generator.HasConsumer,
	}
	spec, err := generator.LoadAsyncApiSpecWithParser(asyncApiSpecPath)
	spec["packageName"] = packageName
	if err != nil {
		return nil, err
	}
	goKafkaGenerator.data = spec

	tmpl := template.Must(template.New("mosaic-kafka-go-code.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-go-code.tmpl"))
	goKafkaGenerator.mosaicTemplate = tmpl

	normalTmpl := template.Must(template.New("kafka-go-code.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/kafka-go-code.tmpl"))
	goKafkaGenerator.normalTemplate = normalTmpl

	mqttTemplate := template.Must(template.New("mqtt-go-code.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mqtt-go-code.tmpl"))
	goKafkaGenerator.mqttTemplate = mqttTemplate
	return &goKafkaGenerator, nil
}

func (thiz *MosaicKafkaGoCodeGenerator) Generate(flavor string) ([]byte, error) {
	var tpl bytes.Buffer
	switch flavor {
	case "mosaic":
		err := thiz.mosaicTemplate.Execute(&tpl, thiz.data)
		if err != nil {
			return nil, err
		}
	case "mqtt":
		err := thiz.mqttTemplate.Execute(&tpl, thiz.data)
		if err != nil {
			return nil, err
		}
	default:
		err := thiz.normalTemplate.Execute(&tpl, thiz.data)
		if err != nil {
			return nil, err
		}
	}

	p, err := format.Source(tpl.Bytes())
	if err != nil {
		return nil, err
	}
	return p, nil
}
