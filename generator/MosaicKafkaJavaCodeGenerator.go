package generator

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
	"os"
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

type MosaicKafkaJavaCodeGenerator struct {
	spec                      *javaSpec
	eventClassTemplate        *template.Template
	producerInterfaceTemplate *template.Template
	consumerInterfaceTemplate *template.Template
}
type javaSpec struct {
	Events   []javaSpecMessage
	Channels map[string]Channel
	Imports  []string
}

type javaSpecMessage struct {
	Name    string
	Message Message
	Typ     string
	Imports []string
}

func NewJavaSpecFromApiSpec(api asyncApiSpec) *javaSpec {
	spec := javaSpec{}
	spec.convertToJavaSpec(api)
	return &spec
}

func (g *javaSpec) convertToJavaSpec(a asyncApiSpec) {
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

func (g *javaSpec) convertAndAddEvent(a asyncApiSpec, value Message, msgType string) {
	var newProps map[string]Property
	imports := []string{}
	if value.Ref != nil {
		value.findMessageByReferenceInComponents(a.Components)
	}
	newProps = a.rewriteProperties(value.Schema.Properties, value.Schema.Required, g.rewriteToJavaProperties)
	value.Schema.Properties = newProps
	for _, prop := range value.Schema.Properties {
		imports = g.checkTypesToImport(prop, imports)
		imports = g.checkAnnotationToImport(prop, imports)
	}

	goMsg := javaSpecMessage{
		Name:    strcase.ToLowerCamel(value.Name),
		Message: value,
		Typ:     msgType,
		Imports: imports,
	}
	g.addEvent(goMsg)
}

func (g *javaSpec) checkAnnotationToImport(prop Property, imports []string) []string {
	if strings.Contains(prop.Type, "@NotNull") {
		toImport := "import javax.validation.constraints.NotNull;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	}
	if strings.Contains(prop.Type, "@NotBlank") {
		toImport := "import javax.validation.constraints.NotBlank;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	}
	if strings.Contains(prop.Type, "@Email") {
		toImport := "import javax.validation.constraints.Email;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	}

	return imports
}

func (g *javaSpec) checkTypesToImport(prop Property, imports []string) []string {
	switch {
	case strings.Contains(prop.Type, "OffsetDateTime"):
		toImport := "import java.time.OffsetDateTime;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	case strings.Contains(prop.Type, "List"):
		toImport := "import java.util.List;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	case strings.Contains(prop.Type, "Map"):
		toImport := "import java.util.Map;"
		if !slices.Contains(imports, toImport) {
			imports = append(imports, toImport)
		}
	}
	return imports
}

func (g *javaSpec) addEvent(message javaSpecMessage) {
	if !containsJavaSpecMessage(g.Events, message) {
		g.Events = append(g.Events, message)
	}
}

func containsJavaSpecMessage(messages []javaSpecMessage, msg javaSpecMessage) bool {
	for _, v := range messages {
		if cmp.Equal(v, msg) {
			return true
		}
	}

	return false
}

func NewMosaicKafkaJavaCodeGenerator(asyncApiSpecPath string) MosaicKafkaJavaCodeGenerator {
	tmpl := template.Must(template.ParseFS(templateFiles, "templates/mosaic-kafka-java-event-class.tmpl"))
	producerInterfaceTmpl := template.Must(template.ParseFS(templateFiles, "templates/mosaic-kafka-java-producer-interface.tmpl"))
	consumerInterfaceTmpl := template.Must(template.ParseFS(templateFiles, "templates/mosaic-kafka-java-consumer-interface.tmpl"))
	spec := loadAsyncApiSpec(asyncApiSpecPath)
	javaSpec := NewJavaSpecFromApiSpec(spec)
	return MosaicKafkaJavaCodeGenerator{
		spec:                      javaSpec,
		eventClassTemplate:        tmpl,
		producerInterfaceTemplate: producerInterfaceTmpl,
		consumerInterfaceTemplate: consumerInterfaceTmpl,
	}
}

func (c MosaicKafkaJavaCodeGenerator) Generate(out string) (string, error) {
	err := os.MkdirAll(out, os.ModePerm)
	if err != nil {
		return "", err
	}
	for _, event := range c.spec.Events {
		s, err := c.createEventClass(out, event)
		if err != nil {
			return s, err
		}
		if event.Typ == "subscribe" {
			_, err = c.createEventProducer(out, event)
			if err != nil {
				return s, err
			}
		}
		if event.Typ == "publish" {
			_, err = c.createEventConsumer(out, event)
			if err != nil {
				return s, err
			}
		}
	}

	return "", nil
}

func (c MosaicKafkaJavaCodeGenerator) createEventClass(out string, event javaSpecMessage) (string, error) {
	var tpl bytes.Buffer
	f, err := os.Create(out + "/" + event.Message.Name + ".java")
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
	return "", nil
}

func (c MosaicKafkaJavaCodeGenerator) createEventProducer(out string, event javaSpecMessage) (string, error) {
	var tpl bytes.Buffer
	f, err := os.Create(out + "/I" + event.Message.Name + "Producer.java")
	if err != nil {
		return "", err
	}
	err = c.producerInterfaceTemplate.Execute(&tpl, event)
	if err != nil {
		return "", err
	}
	_, err = f.Write(tpl.Bytes())
	if err != nil {
		return "", err
	}
	return "", nil
}

func (c MosaicKafkaJavaCodeGenerator) createEventConsumer(out string, event javaSpecMessage) (string, error) {
	var tpl bytes.Buffer
	f, err := os.Create(out + "/I" + event.Message.Name + "Consumer.java")
	if err != nil {
		return "", err
	}
	err = c.consumerInterfaceTemplate.Execute(&tpl, event)
	if err != nil {
		return "", err
	}
	_, err = f.Write(tpl.Bytes())
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *javaSpec) rewriteToJavaProperties(propertyName string, required *[]string, property Property, newProps map[string]Property) {
	fm := "%s %s" //@NotNull type
	typ := ""
	annotation := ""
	if required != nil {
		if slices.Contains(*required, propertyName) {
			if property.Type == "string" {
				if property.Format != nil {
					form := *property.Format
					switch form {
					case "binary":
						annotation = "@NotNull"
					case "email":
						annotation = "@Email"
					default:
						annotation = "@NotBlank"
					}
				}
			} else {
				annotation = "@NotNull"
			}
		}
	}

	switch property.Type {
	case "object":
		if property.Object != nil {
			typ = *property.Object.Name
		} else if property.AdditionalProperties.Type == "string" {
			typ = "Map<String,String>"
		}
	case "array":
		typ = "List<"
		if property.Items.Format != nil {
			typ = typ + typeConversionJavaMap[*property.Items.Format] + ">"
		} else if property.Items.Type == "object" {
			typ = typ + *property.Items.Object.Name + ">"
		} else {
			typ = typ + typeConversionJavaMap[property.Items.Type] + ">"
		}
	default:
		if property.Format != nil {
			typ = typeConversionJavaMap[*property.Format]
		} else {
			typ = typeConversionJavaMap[property.Type]
		}
	}
	wholeString := strings.Trim(fmt.Sprintf(fm, annotation, typ), " ")
	newProps[propertyName] = Property{
		Type:    wholeString,
		Format:  property.Format,
		Minimum: property.Minimum,
		Object:  property.Object,
		Items:   property.Items,
	}
}
