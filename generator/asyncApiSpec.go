package generator

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
	"strings"
)

type asyncApiSpec struct {
	Events   []Message `yaml:"-"`
	AsyncApi string    `yaml:"asyncapi"`
	Info     struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	}
	Servers    map[string]Server
	Channels   map[string]Channel
	Components Components
}

type Server struct {
	Url         string `yaml:"url"`
	Protocol    string `yaml:"protocol"`
	Description string `yaml:"description"`
}

type Channel struct {
	Name      string `yaml:"-"`
	Subscribe *Subscribe
	Publish   *Publish
}

type Subscribe struct {
	Message Message
}

type Publish struct {
	Message Message
}

type Message struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Ref         *string `yaml:"$ref"`
	Schema      Payload `yaml:"payload"`
}

type Payload struct {
	Name                 *string `yaml:"title"`
	Type                 string  `yaml:"type"`
	AdditionalProperties bool    `yaml:"additionalProperties"`
	Properties           map[string]Property
	Ref                  *string  `yaml:"$ref"`
	Required             []string `yaml:"required"`
}

type Property struct {
	Type                 string             `yaml:"type"`
	Format               string             `yaml:"format"`
	Minimum              int                `yaml:"minimum"`
	AdditionalProperties AdditionalProperty `yaml:"additionalProperties"`
	Ref                  *string            `yaml:"$ref"`
	Object               *Payload
	Items                *Item `yaml:"items"`
}

type AdditionalProperty struct {
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

type Components struct {
	Messages map[string]Message
	Schemas  map[string]Payload
}

type Item struct {
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

func (a *asyncApiSpec) convertToUsableStruct() {
	newChannels := make(map[string]Channel)
	for key, value := range a.Channels {
		newKey := strcase.ToCamel(strings.ToLower(key))
		var newProps map[string]Property
		if value.Subscribe != nil {
			if value.Subscribe.Message.Ref != nil {
				value.Subscribe.Message.findMessageByReferenceInComponents(a.Components)
			}
			newProps = a.rewriteProperties(value.Subscribe.Message.Schema.Properties, value.Subscribe.Message.Schema.Required)
			value.Subscribe.Message.Schema.Properties = newProps
			a.addEvent(value.Subscribe.Message)
		}
		if value.Publish != nil {
			if value.Publish.Message.Ref != nil {
				value.Publish.Message.findMessageByReferenceInComponents(a.Components)
			}
			newProps = a.rewriteProperties(value.Publish.Message.Schema.Properties, value.Publish.Message.Schema.Required)
			value.Publish.Message.Schema.Properties = newProps
			a.addEvent(value.Publish.Message)
		}

		value.Name = key
		a.Channels[key] = value

		newChannels[newKey] = a.Channels[key]
	}
	a.Channels = newChannels
}

func (a *asyncApiSpec) rewriteProperties(properties map[string]Property, required []string) map[string]Property {
	newProps := make(map[string]Property)
	props := properties
	for propertyName, property := range props {
		fm := "%s%s `json:\"%s%s\"`" //propertyName optionalPointer type jsonName optionalOmitEmpty
		typ := ""
		pointer := ""
		jsonName := propertyName
		omit := ""
		if !slices.Contains(required, propertyName) {
			pointer = "*"
			omit = ",omitempty"
		}
		switch property.Type {
		case "integer":
			typ = "int"
		case "boolean":
			typ = "bool"
		case "string":
			if property.Format == "date-time" {
				typ = "time.Time"
			} else {
				typ = property.Type
			}

		case "object":
			if property.AdditionalProperties.Type == "string" {
				typ = "map[string]string"
			} else if property.Object != nil {
				typ = "struct {"
				for key, val := range property.Object.Properties {
					typ = typ + "\n\t\t" + strcase.ToCamel(key) + " " + val.Type + "`json:\"" + strcase.ToLowerCamel(key) + "\"`"
				}
				typ = typ + "\n\t}"
			}
		case "array":
			typ = typ + "[]"
			if property.Items.Type == "string" {
				switch property.Items.Format {
				case "binary":
					typ = "[]byte"
				default:
					typ = "string"
				}
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
		}

	}
	return newProps
}

func (a *asyncApiSpec) addEvent(message Message) {
	if !contains(a.Events, message) {
		a.Events = append(a.Events, message)
	}
}

func contains(messages []Message, msg Message) bool {
	for _, v := range messages {
		if cmp.Equal(v, msg) {
			return true
		}
	}

	return false
}

func (m *Message) findMessageByReferenceInComponents(components Components) {
	referenceSlice := strings.Split(*m.Ref, "/")
	messageName := referenceSlice[len(referenceSlice)-1]
	messageFromComponents := components.Messages[messageName]
	m.Description = messageFromComponents.Description
	m.Name = messageFromComponents.Name
	if messageFromComponents.Schema.Ref != nil {
		messageFromComponents.Schema.findPayloadByReferenceInComponents(components)
	}
	m.Schema = messageFromComponents.Schema
}

func (p *Payload) findPayloadByReferenceInComponents(components Components) {
	referenceSlice := strings.Split(*p.Ref, "/")
	payloadName := referenceSlice[len(referenceSlice)-1]
	payloadFromComponents := components.Schemas[payloadName]
	p.AdditionalProperties = payloadFromComponents.AdditionalProperties
	p.Type = payloadFromComponents.Type
	p.Properties = payloadFromComponents.Properties
	p.Required = payloadFromComponents.Required
	for key, prop := range p.Properties {
		if prop.Ref != nil {
			p.findPropertyByReferenceInComponents(key, *prop.Ref, components)
		}
	}
}

func (p *Payload) findPropertyByReferenceInComponents(propertyKey string, propertyRef string, components Components) {
	referenceSlice := strings.Split(propertyRef, "/")
	propertyName := referenceSlice[len(referenceSlice)-1]
	propertyFromComponents := components.Schemas[propertyName]
	//todo ist hier jetzt hart gesetzt...
	prop := Property{
		Type:   "object",
		Object: &propertyFromComponents,
	}
	p.Properties[propertyKey] = prop
}
