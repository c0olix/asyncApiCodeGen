package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goSpec_rewriteToGoProperties(t *testing.T) {
	g := &goSpec{}
	newProps := make(map[string]Property)
	type in struct {
		prop     Property
		required []string
	}

	nestedObjProps := make(map[string]Property)
	nestedObjProps["test"] = Property{
		Type: "string",
	}
	var tests = []struct {
		name string
		in   in
		out  string
	}{
		{
			name: "required integer",
			in: in{
				prop: Property{
					Type: "int32",
				},
				required: []string{"aProperty"},
			},
			out: "int `json:\"aProperty\"`",
		},
		{
			name: "optional integer",
			in: in{
				prop: Property{
					Type: "int32",
				},
				required: []string{},
			},
			out: "*int `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional bool",
			in: in{
				prop: Property{
					Type: "boolean",
				},
				required: []string{},
			},
			out: "*bool `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional string date-time",
			in: in{
				prop: Property{
					Type:   "string",
					Format: "date-time",
				},
				required: []string{},
			},
			out: "*time.Time `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional string",
			in: in{
				prop: Property{
					Type: "string",
				},
				required: []string{},
			},
			out: "*string `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional map string",
			in: in{
				prop: Property{
					Type: "object",
					AdditionalProperties: AdditionalProperty{
						Type: "string",
					},
				},
				required: []string{},
			},
			out: "*map[string]string `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional object",
			in: in{
				prop: Property{
					Type: "object",
					Object: &Payload{
						Name:       strp("NestedObject"),
						Type:       "object",
						Properties: nestedObjProps,
					},
				},

				required: []string{},
			},
			out: "*NestedObject `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional array string",
			in: in{
				prop: Property{
					Type: "array",
					Items: &Item{
						Type: "string",
					},
				},
				required: []string{},
			},
			out: "*[]string `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional array string binary",
			in: in{
				prop: Property{
					Type: "array",
					Items: &Item{
						Type:   "string",
						Format: "binary",
					},
				},
				required: []string{},
			},
			out: "*[][]byte `json:\"aProperty,omitempty\"`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.rewriteToGoProperties("aProperty", tt.in.required, tt.in.prop, newProps)
			assert.Equal(t, tt.out, newProps["AProperty"].Type)
		})
	}
}

func Test_goSpec_convertToGoSpec(t *testing.T) {
	channels := make(map[string]Channel)
	props := make(map[string]Property)
	props["testProp"] = Property{
		Type: "string",
	}
	channels["TEST_CHAN"] = Channel{
		Subscribe: &Subscribe{
			Message: Message{
				Name:        "testChanEvent",
				Description: "a test event",
				Schema: Payload{
					Type:       "object",
					Properties: props,
				},
			},
		},
	}
	a := asyncApiSpec{
		Channels: channels,
	}

	g := NewGoSpecFromApiSpec(a)

	assert.Equal(t, "testChanEvent", g.Channels["TestChan"].Subscribe.Message.Name)
	assert.Equal(t, "a test event", g.Channels["TestChan"].Subscribe.Message.Description)
	assert.Equal(t, "object", g.Channels["TestChan"].Subscribe.Message.Schema.Type)
	assert.Equal(t, "string", g.Channels["TestChan"].Subscribe.Message.Schema.Properties["testProp"].Type)
}
