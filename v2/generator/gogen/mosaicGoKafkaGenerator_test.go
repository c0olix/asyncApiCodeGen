package gogen

import (
	"reflect"
	"testing"
)

func TestMosaicKafkaGoCodeGenerator_getImports(t *testing.T) {
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "successful import on publish",
			args: args{
				data: map[string]interface{}{
					"channels": map[string]interface{}{
						"TEST_CHANNEL": map[string]interface{}{
							"publish": map[string]interface{}{
								"message": map[string]interface{}{
									"payload": map[string]interface{}{
										"properties": map[string]interface{}{
											"testProp": map[string]interface{}{
												"type":   "string",
												"format": "date",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{"time"},
		},
		{
			name: "successful import on subscribe",
			args: args{
				data: map[string]interface{}{
					"channels": map[string]interface{}{
						"TEST_CHANNEL": map[string]interface{}{
							"subscribe": map[string]interface{}{
								"message": map[string]interface{}{
									"payload": map[string]interface{}{
										"properties": map[string]interface{}{
											"testProp": map[string]interface{}{
												"type":   "string",
												"format": "date-time",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{"time"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaGoCodeGenerator{}
			if got := thiz.getImports(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMosaicKafkaGoCodeGenerator_convertToGoType(t *testing.T) {
	type args struct {
		property map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int",
			args: args{
				property: map[string]interface{}{
					"type": "integer",
				},
			},
			want: "int",
		},
		{
			name: "int32",
			args: args{
				property: map[string]interface{}{
					"type":   "intger",
					"format": "int32",
				},
			},
			want: "int32",
		},
		{
			name: "int64",
			args: args{
				property: map[string]interface{}{
					"type": "int64",
				},
			},
			want: "int64",
		},
		{
			name: "string",
			args: args{
				property: map[string]interface{}{
					"type": "string",
				},
			},
			want: "string",
		},
		{
			name: "email",
			args: args{
				property: map[string]interface{}{
					"type":   "string",
					"format": "email",
				},
			},
			want: "string",
		},
		{
			name: "boolean",
			args: args{
				property: map[string]interface{}{
					"type": "boolean",
				},
			},
			want: "bool",
		},
		{
			name: "number",
			args: args{
				property: map[string]interface{}{
					"type": "number",
				},
			},
			want: "float32",
		},
		{
			name: "float",
			args: args{
				property: map[string]interface{}{
					"type":   "number",
					"format": "float",
				},
			},
			want: "float32",
		},
		{
			name: "double",
			args: args{
				property: map[string]interface{}{
					"type":   "number",
					"format": "double",
				},
			},
			want: "float64",
		},
		{
			name: "binary",
			args: args{
				property: map[string]interface{}{
					"type":   "string",
					"format": "binary",
				},
			},
			want: "[]byte",
		},
		{
			name: "date",
			args: args{
				property: map[string]interface{}{
					"type": "date",
				},
			},
			want: "time.Time",
		},
		{
			name: "date-time",
			args: args{
				property: map[string]interface{}{
					"type": "date-time",
				},
			},
			want: "time.Time",
		},
		{
			name: "map string string",
			args: args{
				property: map[string]interface{}{
					"type": "object",
					"additionalProperties": map[string]interface{}{
						"type": "string",
					},
				},
			},
			want: "map[string]string",
		},
		{
			name: "object",
			args: args{
				property: map[string]interface{}{
					"type":                 "object",
					"title":                "TestObject",
					"additionalProperties": false,
				},
			},
			want: "TestObject",
		},
		{
			name: "object",
			args: args{
				property: map[string]interface{}{
					"type":  "object",
					"title": "TestObject",
				},
			},
			want: "TestObject",
		},
		{
			name: "array of objects",
			args: args{
				property: map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type":  "object",
						"title": "TestObject",
					},
				},
			},
			want: "[]TestObject",
		},
		{
			name: "array of files",
			args: args{
				property: map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type":   "string",
						"format": "binary",
					},
				},
			},
			want: "[][]byte",
		},
		{
			name: "array of strings",
			args: args{
				property: map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			want: "[]string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaGoCodeGenerator{}
			if got := thiz.convertToGoType(tt.args.property); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToGoType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMosaicKafkaGoCodeGenerator_validations(t *testing.T) {

	type args struct {
		property map[string]interface{}
		required bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "required email",
			args: args{
				property: map[string]interface{}{
					"format": "email",
				},
				required: true,
			},
			want: ` validate:"required,email"`,
		},
		{
			name: "int min",
			args: args{
				property: map[string]interface{}{
					"minimum": 1,
					"type":    "integer",
				},
				required: false,
			},
			want: ` validate:"gte=1"`,
		},
		{
			name: "int exclusive min",
			args: args{
				property: map[string]interface{}{
					"exclusiveMinimum": 1,
					"type":             "integer",
				},
				required: false,
			},
			want: ` validate:"gt=1"`,
		},
		{
			name: "int max",
			args: args{
				property: map[string]interface{}{
					"maximum": 1,
					"type":    "integer",
				},
				required: false,
			},
			want: ` validate:"lte=1"`,
		},
		{
			name: "int exclusive max",
			args: args{
				property: map[string]interface{}{
					"exclusiveMaximum": 1,
					"type":             "integer",
				},
				required: false,
			},
			want: ` validate:"lt=1"`,
		},
		{
			name: "string min",
			args: args{
				property: map[string]interface{}{
					"minLength": 1,
					"type":      "string",
				},
				required: false,
			},
			want: ` validate:"min=1"`,
		},
		{
			name: "string max",
			args: args{
				property: map[string]interface{}{
					"maxLength": 1,
					"type":      "string",
				},
				required: false,
			},
			want: ` validate:"max=1"`,
		},
		{
			name: "array min",
			args: args{
				property: map[string]interface{}{
					"minItems": 1,
					"type":     "array",
				},
				required: false,
			},
			want: ` validate:"min=1"`,
		},
		{
			name: "array max",
			args: args{
				property: map[string]interface{}{
					"maxItems": 1,
					"type":     "array",
				},
				required: false,
			},
			want: ` validate:"max=1"`,
		},
		{
			name: "array unique",
			args: args{
				property: map[string]interface{}{
					"uniqueItems": true,
					"type":        "array",
				},
				required: false,
			},
			want: ` validate:"unique"`,
		},
		{
			name: "not required",
			args: args{
				property: map[string]interface{}{
					"type": "string",
				},
				required: false,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaGoCodeGenerator{}
			if got := thiz.validations(tt.args.property, tt.args.required); got != tt.want {
				t.Errorf("validations() = %v, want %v", got, tt.want)
			}
		})
	}
}
