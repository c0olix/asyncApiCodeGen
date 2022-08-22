package generator

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"
	"go/format"
	"golang.org/x/exp/slices"
	"os"
	"strings"
	"text/template"
)

type goSpec struct {
	Events   []goSpecMessage
	Imports  []string
	Channels map[string]Channel
}

type goSpecMessage struct {
	Message Message
	Typ     string
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
			g.convertAndAddEvent(a, value.Subscribe.Message, "subscribe")
		} else if value.Publish != nil {
			g.convertAndAddEvent(a, value.Publish.Message, "publish")
		}

		value.Name = key
		a.Channels[key] = value

		newChannels[newKey] = a.Channels[key]
	}
	g.Channels = newChannels
}

func (g *goSpec) convertAndAddEvent(a asyncApiSpec, value Message, msgType string) {
	var newProps map[string]Property
	if value.Ref != nil {
		value.findMessageByReferenceInComponents(a.Components)
	}
	newProps = a.rewriteProperties(value.Schema.Properties, value.Schema.Required, g.rewriteToGoProperties)
	value.Schema.Properties = newProps
	goMsg := goSpecMessage{
		Message: value,
		Typ:     msgType,
	}
	g.addEvent(goMsg)
}

func (g *goSpec) addEvent(message goSpecMessage) {
	if !contains(g.Events, message) {
		g.Events = append(g.Events, message)
	}
}

func (g *goSpec) rewriteToGoProperties(propertyName string, required []string, property Property, newProps map[string]Property) {
	fm := "%s%s `json:\"%s%s\"`" //optionalPointer type jsonName optionalOmitEmpty
	typ := ""
	pointer := ""
	jsonName := propertyName
	omit := ""
	if !slices.Contains(required, propertyName) {
		pointer = "*"
		omit = ",omitempty"
	}
	switch property.Type {
	case "int32":
		typ = "int"
	case "int64":
		typ = "int64"
	case "boolean":
		typ = "bool"
	case "string":
		if property.Format == "date-time" {
			typ = "time.Time"
			g.Imports = append(g.Imports, "time")
		} else {
			typ = property.Type
		}

	case "object":
		if property.AdditionalProperties.Type == "string" {
			typ = "map[string]string"
		} else if property.Object != nil {
			typ = *property.Object.Name
		}
	case "array":
		typ = "[]"
		switch property.Items.Type {
		case "string":
			switch property.Items.Format {
			case "binary":
				typ = typ + "[]byte"
			default:
				typ = typ + "string"
			}
		case "object":
			typ = typ + *property.Items.Object.Name
		}

	default:
		typ = property.Type
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

func (c MosaicKafkaGoCodeGenerator) Generate(out string) (string, error) {
	var tpl bytes.Buffer
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	err = c.template.Execute(&tpl, c.spec)
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
