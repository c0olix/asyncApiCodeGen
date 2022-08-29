package javagen

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

var typeConversionJavaMap = map[string]string{
	"int32":     "Integer",
	"int64":     "Long",
	"string":    "String",
	"email":     "String",
	"boolean":   "Boolean",
	"float":     "Float",
	"double":    "Double",
	"binary":    "File",
	"date":      "OffsetDateTime",
	"date-time": "OffsetDateTime",
	"password":  "password",
}

//go:embed templates
var templateFiles embed.FS

type MosaicKafkaJavaCodeGenerator struct {
	eventClassTemplate        *template.Template
	producerInterfaceTemplate *template.Template
	consumerInterfaceTemplate *template.Template
	data                      map[string]interface{}
	log                       *logrus.Logger
}

type MosaicKafkaJavaCodeResult struct {
	Files []ResultFile
}

type ResultFile struct {
	Name    string
	Content []byte
}

func (c MosaicKafkaJavaCodeGenerator) getImports(messagePayload map[string]interface{}) []string {
	var out []string
	properties := messagePayload["properties"].(map[string]interface{})
	var required []interface{}
	if messagePayload["required"] != nil {
		required = messagePayload["required"].([]interface{})
	}
	for propKey, prop := range properties {

		property := prop.(map[string]interface{})
		typ := property["type"].(string)
		var format *string
		if property["format"] != nil {
			form := property["format"].(string)
			format = &form
		}
		if required != nil {
			for _, reqProp := range required {
				if propKey == reqProp.(string) {
					if typ == "string" && format == nil {
						importStatement := "import javax.validation.constraints.NotBlank;"
						if !slices.Contains(out, importStatement) {
							out = append(out, importStatement)
						}
					} else {
						importStatement := "import javax.validation.constraints.NotNull;"
						if !slices.Contains(out, importStatement) {
							out = append(out, importStatement)
						}
					}
				}
			}
		}
		if format != nil {
			switch *format {
			case "email":
				importStatement := "import javax.validation.constraints.Email;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			case "binary":
				importStatement := "import java.io.File;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			}
		} else {
			switch typ {
			case "date", "date-time":
				importStatement := "import java.time.OffsetDateTime"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			case "array":
				importStatement := "import java.util.List;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			case "object":
				if property["additionalProperties"] != nil {
					_, ok := property["additionalProperties"].(map[string]interface{})
					if ok {
						importStatement := "import java.util.Map;"
						if !slices.Contains(out, importStatement) {
							out = append(out, importStatement)
						}
					}
					additionalProperties, ok := property["additionalProperties"].(bool)
					if ok && additionalProperties {
						importStatement := "import java.util.Map;"
						if !slices.Contains(out, importStatement) {
							out = append(out, importStatement)
						}
					}
				}
			}
		}
	}
	return out
}

func (c MosaicKafkaJavaCodeGenerator) convertToJavaType(property map[string]interface{}) string {
	switch property["type"] {
	case "object":
		if property["additionalProperties"] != nil {
			additionalProperties, ok := property["additionalProperties"].(map[string]interface{})
			if ok {
				if additionalProperties["type"] == "string" {
					return "Map<String,String>"
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
		typ := "List<"
		if property["items"] != nil {
			items := property["items"].(map[string]interface{})
			if items["format"] != nil {
				typ = typ + typeConversionJavaMap[items["format"].(string)] + ">"
			} else if items["type"] == "object" {
				return items["title"].(string)
			} else {
				typ = typ + typeConversionJavaMap[items["type"].(string)] + ">"
			}
			return typ
		}
	default:
		if property["format"] != nil {
			return typeConversionJavaMap[property["format"].(string)]
		} else {
			return typeConversionJavaMap[property["type"].(string)]
		}
	}
	return ""
}

func (c MosaicKafkaJavaCodeGenerator) getAnnotations(propertyName string, property map[string]interface{}, required []interface{}) []string {
	var annotations []string
	if generator.HasDefault(property) {
		annotations = append(annotations, "@Builder.Default")
	}
	if generator.CheckRequired(propertyName, required) {
		if property["type"] == "string" && property["format"] == nil {
			annotations = append(annotations, "@NotBlank")
		} else {
			annotations = append(annotations, "@NotNull")
		}
	}
	if property["format"] != nil {
		form := property["format"].(string)
		if form == "email" {
			annotations = append(annotations, "@Email")
		}
	}
	return annotations
}

func NewMosaicKafkaJavaCodeGenerator(asyncApiSpecPath string, logger *logrus.Logger) (*MosaicKafkaJavaCodeGenerator, error) {
	javaKafkaGenerator := MosaicKafkaJavaCodeGenerator{
		log: logger,
	}
	fns := template.FuncMap{
		"getImports":        javaKafkaGenerator.getImports,
		"getMessages":       generator.GetMessages,
		"getObjects":        generator.GetNestedObjects,
		"getItemObjects":    generator.GetItemObjects,
		"hasDefault":        generator.HasDefault,
		"getAnnotations":    javaKafkaGenerator.getAnnotations,
		"convertToJavaType": javaKafkaGenerator.convertToJavaType,
		"lower":             strings.ToLower,
		"camel":             strcase.ToCamel,
		"checkRequired":     generator.CheckRequired,
	}
	spec, err := generator.LoadAsyncApiSpecWithParser(asyncApiSpecPath)
	if err != nil {
		return nil, err
	}
	javaKafkaGenerator.data = spec

	tmpl := template.Must(template.New("mosaic-kafka-java-event-class.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-event-class.tmpl"))
	javaKafkaGenerator.eventClassTemplate = tmpl
	tmplProducer := template.Must(template.New("mosaic-kafka-java-producer-interface.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-producer-interface.tmpl"))
	javaKafkaGenerator.producerInterfaceTemplate = tmplProducer
	tmplConsumer := template.Must(template.New("mosaic-kafka-java-consumer-interface.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-consumer-interface.tmpl"))
	javaKafkaGenerator.consumerInterfaceTemplate = tmplConsumer
	return &javaKafkaGenerator, nil
}

func (c MosaicKafkaJavaCodeGenerator) Generate() (*MosaicKafkaJavaCodeResult, error) {
	var results []ResultFile
	messages := generator.GetMessages(c.data)
	for _, message := range messages {
		var tpl bytes.Buffer
		err := c.eventClassTemplate.Execute(&tpl, message)
		if err != nil {
			return nil, err
		}
		results = append(results, ResultFile{
			Name:    message["name"].(string),
			Content: tpl.Bytes(),
		})
	}
	for _, channel := range c.data["channels"].(map[string]interface{}) {
		var tpl bytes.Buffer
		ch := channel.(map[string]interface{})
		if ch["publish"] != nil {
			operation := ch["publish"].(map[string]interface{})
			err := c.consumerInterfaceTemplate.Execute(&tpl, operation)
			if err != nil {
				return nil, err
			}
			message := operation["message"].(map[string]interface{})
			results = append(results, ResultFile{
				Name:    "I" + message["name"].(string) + "Consumer",
				Content: tpl.Bytes(),
			})
		} else {
			operation := ch["subscribe"].(map[string]interface{})
			err := c.producerInterfaceTemplate.Execute(&tpl, operation)
			if err != nil {
				return nil, err
			}
			message := operation["message"].(map[string]interface{})
			results = append(results, ResultFile{
				Name:    "I" + message["name"].(string) + "Producer",
				Content: tpl.Bytes(),
			})
		}
	}
	return &MosaicKafkaJavaCodeResult{
		Files: results,
	}, nil
}
