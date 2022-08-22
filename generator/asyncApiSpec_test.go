package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_findReferenceInComponents(t *testing.T) {
	messages := make(map[string]Message)
	messages["TestEvent"] = Message{
		Name:        "TestEvent",
		Description: "A Test Event",
		Ref:         nil,
		Schema:      Payload{},
	}
	comps := Components{
		Messages: messages,
		Schemas:  nil,
	}
	m := Message{
		Ref: strp("#/components/messages/TestEvent"),
	}
	m.findMessageByReferenceInComponents(comps)

	assert.Equal(t, "TestEvent", m.Name)
	assert.Equal(t, "A Test Event", m.Description)
	assert.Equal(t, strp("#/components/messages/TestEvent"), m.Ref)
	assert.Equal(t, Payload{}, m.Schema)
}

func TestPayload_findReferenceInComponents(t *testing.T) {
	payloads := make(map[string]Payload)
	properties := make(map[string]Property)
	properties["email"] = Property{
		Type:    "string",
		Format:  strp("email"),
		Minimum: ip(0),
	}
	payloads["TestPayload"] = Payload{
		Type:                 "object",
		AdditionalProperties: bp(false),
		Properties:           properties,
		Ref:                  nil,
	}
	comps := Components{
		Messages: nil,
		Schemas:  payloads,
	}
	p := Payload{
		Ref: strp("#/components/schemas/TestPayload"),
	}
	p.findPayloadByReferenceInComponents(comps)

	assert.Equal(t, "object", p.Type)
	assert.Equal(t, properties, p.Properties)
	assert.Equal(t, strp("#/components/schemas/TestPayload"), p.Ref)
	assert.Equal(t, bp(false), p.AdditionalProperties)
}

func TestPayload_findReferenceInComponentsWithPorpertyWithRef(t *testing.T) {
	payloads := make(map[string]Payload)
	properties := make(map[string]Property)
	properties["AnObject"] = Property{
		Ref: strp("#/components/schemas/AnObject"),
	}
	payloads["TestPayload"] = Payload{
		Type:                 "object",
		AdditionalProperties: bp(false),
		Properties:           properties,
		Ref:                  nil,
	}
	propertiesForObject := make(map[string]Property)
	propertiesForObject["AField"] = Property{
		Type: "string",
	}
	payloads["AnObject"] = Payload{
		Type:                 "object",
		AdditionalProperties: bp(false),
		Properties:           propertiesForObject,
		Ref:                  nil,
	}
	comps := Components{
		Messages: nil,
		Schemas:  payloads,
	}
	p := Payload{
		Ref: strp("#/components/schemas/TestPayload"),
	}
	p.findPayloadByReferenceInComponents(comps)

	assert.Equal(t, "object", p.Type)
	assert.Equal(t, properties, p.Properties)
	assert.Equal(t, strp("#/components/schemas/TestPayload"), p.Ref)
	assert.Equal(t, bp(false), p.AdditionalProperties)
	assert.Equal(t, "string", p.Properties["AnObject"].Object.Properties["AField"].Type)
}

func TestRewriteProperties(t *testing.T) {
	a := asyncApiSpec{}

	properties := make(map[string]Property)
	properties["anInteger"] = Property{
		Type: "integer",
	}
	properties["aBoolean"] = Property{
		Type: "boolean",
	}
	required := []string{
		"aBoolean",
	}

	propertyRewriteFunc := func(propertyName string, required *[]string, property Property, newProps map[string]Property) {
		if propertyName == "anInteger" {
			newProps[propertyName] = Property{
				Type: "*int",
			}
		} else {
			newProps[propertyName] = Property{
				Type: "bool",
			}
		}
	}

	newProps := a.rewriteProperties(properties, &required, propertyRewriteFunc)
	assert.Equal(t, "*int", newProps["anInteger"].Type)
	assert.Equal(t, "bool", newProps["aBoolean"].Type)
}

func bp(in bool) *bool {
	return &in
}

func ip(in int) *int {
	return &in
}
