package generator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_goSpec_rewriteToGoProperties(t *testing.T) {
	g := &goSpec{}
	newProps := make(map[string]Property)
	type in struct {
		prop     Property
		required *[]string
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
				required: &[]string{"aProperty"},
			},
			out: "int `json:\"aProperty\"`",
		},
		{
			name: "optional integer",
			in: in{
				prop: Property{
					Type: "int32",
				},
				required: nil,
			},
			out: "*int `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional bool",
			in: in{
				prop: Property{
					Type: "boolean",
				},
				required: nil,
			},
			out: "*bool `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional string date-time",
			in: in{
				prop: Property{
					Type:   "string",
					Format: strp("date-time"),
				},
				required: nil,
			},
			out: "*time.Time `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional string",
			in: in{
				prop: Property{
					Type: "string",
				},
				required: nil,
			},
			out: "*string `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional map string",
			in: in{
				prop: Property{
					Type: "object",
					AdditionalProperties: &AdditionalProperty{
						Type: "string",
					},
				},
				required: nil,
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

				required: nil,
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
				required: nil,
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
						Format: strp("binary"),
					},
				},
				required: nil,
			},
			out: "*[][]byte `json:\"aProperty,omitempty\"`",
		},
		{
			name: "optional array object",
			in: in{
				prop: Property{
					Type: "array",
					Items: &Item{
						Type: "object",
						Object: &Payload{
							Name: strp("TestItem"),
							Type: "string",
						},
					},
				},
				required: nil,
			},
			out: "*[]TestItem `json:\"aProperty,omitempty\"`",
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

func TestMosaicKafkaGoCodeGenerator_Generate(t *testing.T) {
	expected, err := os.ReadFile("./test-spec/expected/out.gen.go")
	assert.Nil(t, err)

	gen := NewMosaicKafkaGoCodeGenerator("./test-spec/test-spec.yaml")
	out, err := gen.Generate()
	assert.Nil(t, err)
	assert.Equal(t, string(expected), string(out))
}
