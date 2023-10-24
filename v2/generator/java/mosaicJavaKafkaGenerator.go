package javagen

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/c0olix/asyncApiCodeGen/v2/generator"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
	"text/template"
)

var typeConversionJavaMap = map[string]string{
	"integer":   "Integer",
	"int32":     "Integer",
	"int64":     "Long",
	"string":    "String",
	"email":     "String",
	"boolean":   "Boolean",
	"number":    "Float",
	"float":     "Float",
	"double":    "Double",
	"binary":    "File",
	"date":      "OffsetDateTime",
	"date-time": "OffsetDateTime",
}

var flavor = "springboot2"

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

func (thiz *MosaicKafkaJavaCodeGenerator) getImports(messagePayload map[string]interface{}) []string {
	var out []string
	var validationPackage = "javax"
	if flavor == "springboot3" {
		validationPackage = "jakarta"
	}
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
		if property["minimum"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Min;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if property["maximum"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Max;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if property["minLength"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Size;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if property["maxLength"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Size;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if property["minItems"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Size;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if property["maxItems"] != nil {
			importStatement := fmt.Sprintf("import %s.validation.constraints.Size;", validationPackage)
			if !slices.Contains(out, importStatement) {
				out = append(out, importStatement)
			}
		}
		if required != nil {
			for _, reqProp := range required {
				if propKey == reqProp.(string) {
					if typ == "string" && format == nil {
						importStatement := fmt.Sprintf("import %s.validation.constraints.NotBlank;", validationPackage)
						if !slices.Contains(out, importStatement) {
							out = append(out, importStatement)
						}
					} else {
						importStatement := fmt.Sprintf("import %s.validation.constraints.NotNull;", validationPackage)
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
				importStatement := fmt.Sprintf("import %s.validation.constraints.Email;", validationPackage)
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			case "binary":
				importStatement := "import java.io.File;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			case "date", "date-time":
				importStatement := "import java.time.OffsetDateTime;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
			}
		} else {
			switch typ {
			case "array":
				importStatement := "import java.util.List;"
				if !slices.Contains(out, importStatement) {
					out = append(out, importStatement)
				}
				prop := prop.(map[string]interface{})
				if prop["items"] != nil {
					items := prop["items"].(map[string]interface{})
					if items["format"] != nil {
						switch items["format"] {
						case "email":
							importStatement := fmt.Sprintf("import %s.validation.constraints.Email;", validationPackage)
							if !slices.Contains(out, importStatement) {
								out = append(out, importStatement)
							}
						case "binary":
							importStatement := "import java.io.File;"
							if !slices.Contains(out, importStatement) {
								out = append(out, importStatement)
							}
						case "date", "date-time":
							importStatement := "import java.time.OffsetDateTime;"
							if !slices.Contains(out, importStatement) {
								out = append(out, importStatement)
							}
						}
					}
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
	sort.Strings(out)
	return out
}

func (thiz *MosaicKafkaJavaCodeGenerator) convertToJavaType(property map[string]interface{}) string {
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
				typ = typ + items["title"].(string) + ">"
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

func (thiz *MosaicKafkaJavaCodeGenerator) getAnnotations(propertyName string, property map[string]interface{}, required []interface{}) []string {
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
	if property["minimum"] != nil {
		min := property["minimum"].(int)
		annotations = append(annotations, fmt.Sprintf("@Min(%d)", min))
	}
	if property["maximum"] != nil {
		max := property["maximum"].(int)
		annotations = append(annotations, fmt.Sprintf("@Max(%d)", max))
	}
	maxLengthAlreadyProcessed := false
	if property["minLength"] != nil {
		min := property["minLength"].(int)
		if property["maxLength"] != nil {
			max := property["maxLength"].(int)
			annotations = append(annotations, fmt.Sprintf("@Size(min=%d,max=%d)", min, max))
			maxLengthAlreadyProcessed = true
		} else {
			annotations = append(annotations, fmt.Sprintf("@Size(min=%d)", min))
		}
	}
	if property["maxLength"] != nil && !maxLengthAlreadyProcessed {
		max := property["maxLength"].(int)
		annotations = append(annotations, fmt.Sprintf("@Size(max=%d)", max))
	}
	maxItemsAlreadyProcessed := false
	if property["minItems"] != nil {
		min := property["minItems"].(int)
		if property["maxItems"] != nil {
			max := property["maxItems"].(int)
			annotations = append(annotations, fmt.Sprintf("@Size(min=%d,max=%d)", min, max))
			maxItemsAlreadyProcessed = true
		} else {
			annotations = append(annotations, fmt.Sprintf("@Size(min=%d)", min))
		}
	}
	if property["maxItems"] != nil && !maxItemsAlreadyProcessed {
		max := property["maxItems"].(int)
		annotations = append(annotations, fmt.Sprintf("@Size(max=%d)", max))
	}
	sort.Strings(annotations)
	return annotations
}

func NewMosaicKafkaJavaCodeGenerator(asyncApiSpecPath string, packageName string, fl string, logger *logrus.Logger) (*MosaicKafkaJavaCodeGenerator, error) {
	javaKafkaGenerator := MosaicKafkaJavaCodeGenerator{
		log: logger,
	}
	if fl != "" {
		flavor = fl
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
	spec["packageName"] = packageName
	javaKafkaGenerator.data = spec

	tmpl := template.Must(template.New("mosaic-kafka-java-event-class.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-event-class.tmpl"))
	javaKafkaGenerator.eventClassTemplate = tmpl
	tmplProducer := template.Must(template.New("mosaic-kafka-java-producer-interface.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-producer-interface.tmpl"))
	javaKafkaGenerator.producerInterfaceTemplate = tmplProducer
	tmplConsumer := template.Must(template.New("mosaic-kafka-java-consumer-interface.tmpl").Funcs(fns).ParseFS(templateFiles, "templates/mosaic-kafka-java-consumer-interface.tmpl"))
	javaKafkaGenerator.consumerInterfaceTemplate = tmplConsumer
	return &javaKafkaGenerator, nil
}

func (thiz *MosaicKafkaJavaCodeGenerator) Generate() (*MosaicKafkaJavaCodeResult, error) {
	var results []ResultFile
	messages := generator.GetMessages(thiz.data)
	for _, message := range messages {
		var tpl bytes.Buffer
		message["packageName"] = thiz.data["packageName"]
		err := thiz.eventClassTemplate.Execute(&tpl, message)
		if err != nil {
			return nil, err
		}
		results = append(results, ResultFile{
			Name:    message["name"].(string),
			Content: tpl.Bytes(),
		})
	}
	objects := generator.GetNestedObjects(messages)
	for _, obj := range objects {
		var tpl bytes.Buffer
		objEvent := map[string]interface{}{
			"name": obj["title"],
			"payload": map[string]interface{}{
				"additionalProperties": obj["additionalProperties"],
				"title":                obj["title"],
				"required":             obj["required"],
				"properties":           obj["properties"],
				"type":                 obj["type"],
				"parent":               obj["parent"],
			},
			"packageName": thiz.data["packageName"],
		}
		err := thiz.eventClassTemplate.Execute(&tpl, objEvent)
		if err != nil {
			return nil, err
		}
		results = append(results, ResultFile{
			Name:    objEvent["name"].(string),
			Content: tpl.Bytes(),
		})
	}
	items := generator.GetItemObjects(messages)
	for _, obj := range items {
		var tpl bytes.Buffer
		items := obj["items"].(map[string]interface{})
		objEvent := map[string]interface{}{
			"name": items["title"],
			"payload": map[string]interface{}{
				"additionalProperties": items["additionalProperties"],
				"title":                items["title"],
				"required":             items["required"],
				"properties":           items["properties"],
				"type":                 items["type"],
				"parent":               items["parent"],
			},
			"packageName": thiz.data["packageName"],
		}
		err := thiz.eventClassTemplate.Execute(&tpl, objEvent)
		if err != nil {
			return nil, err
		}
		results = append(results, ResultFile{
			Name:    objEvent["name"].(string),
			Content: tpl.Bytes(),
		})
	}
	for _, channel := range thiz.data["channels"].(map[string]interface{}) {
		var tpl bytes.Buffer
		ch := channel.(map[string]interface{})
		if ch["publish"] != nil {
			operation := ch["publish"].(map[string]interface{})
			operation["packageName"] = thiz.data["packageName"]
			err := thiz.consumerInterfaceTemplate.Execute(&tpl, operation)
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
			operation["packageName"] = thiz.data["packageName"]
			if flavor == "springboot3" {
				operation["validationPackage"] = "jakarta"
			} else {
				operation["validationPackage"] = "javax"
			}
			err := thiz.producerInterfaceTemplate.Execute(&tpl, operation)
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
