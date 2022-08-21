package generator

import (
	"strings"
)

type propertyRewriteFunc func(propertyName string, required []string, property Property, newProps map[string]Property)

type asyncApiSpec struct {
	AsyncApi string `yaml:"asyncapi"`
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

func (a *asyncApiSpec) rewriteProperties(properties map[string]Property, required []string, conversationFunc propertyRewriteFunc) map[string]Property {
	newProps := make(map[string]Property)
	props := properties
	for propertyName, property := range props {
		if property.Object != nil {
			objProps := a.rewriteProperties(property.Object.Properties, property.Object.Required, conversationFunc)
			property.Object.Properties = objProps
		}
		conversationFunc(propertyName, required, property, newProps)
	}
	return newProps
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
	prop := p.Properties[propertyKey]
	prop.Object = &propertyFromComponents
	prop.Type = propertyFromComponents.Type
	p.Properties[propertyKey] = prop
}
