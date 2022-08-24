package generator

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"
	"go/format"
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

type goSpec struct {
	Events   []goSpecMessage
	Imports  []string
	Channels map[string]Channel
}

type goSpecMessage struct {
	Message     Message
	Typ         string
	OperationId string
	Topic       string
}

func NewGoSpecFromApiSpec(api asyncApiSpec) *goSpec {
	spec := goSpec{}
	spec.convertToGoSpec(api)
	return &spec
}

func (g *goSpec) convertToGoSpec(a asyncApiSpec) {
	newChannels := make(map[string]Channel)
	for key, value := range a.Channels {
		newKey := strcase.ToCamel(strings.ToLower(key))

		if value.Subscribe != nil {
			if value.Subscribe.OperationId == nil {
				opId := "Produce" + value.Subscribe.Message.Name
				value.Subscribe.OperationId = &opId
			} else {
				opId := strcase.ToCamel(*value.Subscribe.OperationId)
				value.Subscribe.OperationId = &opId
			}
			g.convertAndAddEvent(a, value.Subscribe.Message, "subscribe", *value.Subscribe.OperationId, newKey)
		} else if value.Publish != nil {
			if value.Publish.OperationId == nil {
				opId := "Consume" + value.Publish.Message.Name
				value.Publish.OperationId = &opId
			} else {
				opId := strcase.ToCamel(*value.Publish.OperationId)
				value.Publish.OperationId = &opId
			}
			g.convertAndAddEvent(a, value.Publish.Message, "publish", *value.Publish.OperationId, newKey)
		}

		value.Name = key
		a.Channels[key] = value

		newChannels[newKey] = a.Channels[key]
	}
	g.Channels = newChannels
}

func (g *goSpec) convertAndAddEvent(a asyncApiSpec, value Message, msgType string, opId string, topicName string) {
	var newProps map[string]Property
	if value.Ref != nil {
		value.findMessageByReferenceInComponents(a.Components)
	}
	newProps = a.rewriteProperties(value.Schema.Properties, value.Schema.Required, g.rewriteToGoProperties)
	value.Schema.Properties = newProps
	goMsg := goSpecMessage{
		Message:     value,
		Typ:         msgType,
		OperationId: opId,
		Topic:       topicName,
	}
	g.addEvent(goMsg)
}

func (g *goSpec) addEvent(message goSpecMessage) {
	if !contains(g.Events, message) {
		g.Events = append(g.Events, message)
	}
}

func (g *goSpec) rewriteToGoProperties(propertyName string, required *[]string, property Property, newProps map[string]Property) {
	fm := "%s%s `json:\"%s%s\"`" //optionalPointer type jsonName optionalOmitEmpty
	typ := ""
	pointer := ""
	jsonName := propertyName
	omit := ""
	if required != nil {
		if !slices.Contains(*required, propertyName) {
			pointer = "*"
			omit = ",omitempty"
		}
	} else {
		pointer = "*"
		omit = ",omitempty"
	}
	switch property.Type {
	case "object":
		if property.Object != nil {
			typ = *property.Object.Name
		} else if property.AdditionalProperties.Type == "string" {
			typ = "map[string]string"
		}
	case "array":
		typ = "[]"
		if property.Items.Format != nil {
			typ = typ + typeConversionGoMap[*property.Items.Format]
			if strings.Contains(typ, "date") {
				g.Imports = append(g.Imports, "time")
			}
		} else if property.Items.Type == "object" {
			typ = typ + *property.Items.Object.Name
		} else {
			typ = typ + typeConversionGoMap[property.Items.Type]
		}
	default:
		if property.Format != nil {
			typ = typeConversionGoMap[*property.Format]
			if strings.Contains(typ, "time") {
				g.Imports = append(g.Imports, "time")
			}
		} else {
			typ = typeConversionGoMap[property.Type]
		}
	}
	wholeString := fmt.Sprintf(fm, pointer, typ, jsonName, omit)
	newPropertyName := strcase.ToCamel(propertyName)
	newProps[newPropertyName] = Property{
		Type:    wholeString,
		Format:  property.Format,
		Minimum: property.Minimum,
		Object:  property.Object,
		Items:   property.Items,
	}
}

func contains(messages []goSpecMessage, msg goSpecMessage) bool {
	for _, v := range messages {
		if cmp.Equal(v, msg) {
			return true
		}
	}

	return false
}

type MosaicKafkaGoCodeGenerator struct {
	template *template.Template
	spec     *goSpec
}

func NewMosaicKafkaGoCodeGenerator(asyncApiSpecPath string) MosaicKafkaGoCodeGenerator {
	spec := loadAsyncApiSpec(asyncApiSpecPath)
	goSpec := NewGoSpecFromApiSpec(spec)

	tmpl := template.Must(template.ParseFS(templateFiles, "templates/mosaic-kafka-go-code.tmpl"))
	return MosaicKafkaGoCodeGenerator{
		template: tmpl,
		spec:     goSpec,
	}
}

func (c MosaicKafkaGoCodeGenerator) Generate() ([]byte, error) {
	var tpl bytes.Buffer
	err := c.template.Execute(&tpl, c.spec)
	if err != nil {
		return nil, err
	}
	p, err := format.Source(tpl.Bytes())
	if err != nil {
		return nil, err
	}
	return p, nil
}
