package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_javaSpec_rewriteToJavaProperties(t *testing.T) {
	g := &javaSpec{}
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
			out: "@NotNull Integer",
		},
		{
			name: "optional integer",
			in: in{
				prop: Property{
					Type: "int32",
				},
				required: nil,
			},
			out: " Integer",
		},
		{
			name: "optional bool",
			in: in{
				prop: Property{
					Type: "boolean",
				},
				required: nil,
			},
			out: " Boolean",
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
			out: " OffsetDateTime",
		},
		{
			name: "optional string",
			in: in{
				prop: Property{
					Type: "string",
				},
				required: nil,
			},
			out: " String",
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
			out: " Map<String,String>",
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
			out: " NestedObject",
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
			out: " List<String>",
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
			out: " List<File>",
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
			out: " List<TestItem>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.rewriteToJavaProperties("aProperty", tt.in.required, tt.in.prop, newProps)
			assert.Equal(t, tt.out, newProps["aProperty"].Type)
		})
	}
}

func Test_javaSpec_convertToJavaSpec(t *testing.T) {
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
