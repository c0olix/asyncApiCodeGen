package javagen

import (
	"reflect"
	"sort"
	"testing"
)

func TestMosaicKafkaJavaCodeGenerator_getImports(t *testing.T) {
	type args struct {
		messagePayload map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "all imports",
			args: args{
				messagePayload: map[string]interface{}{
					"required": []interface{}{
						"time",
						"string",
					},
					"properties": map[string]interface{}{
						"time": map[string]interface{}{
							"type":   "string",
							"format": "date-time",
						},
						"binary": map[string]interface{}{
							"type":   "string",
							"format": "binary",
						},
						"email": map[string]interface{}{
							"type":   "string",
							"format": "email",
						},
						"string": map[string]interface{}{
							"type": "string",
						},
						"array": map[string]interface{}{
							"type": "array",
						},
						"map": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": map[string]interface{}{},
						},
						"map true": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
						},
					},
				},
			},
			want: []string{
				"import java.time.OffsetDateTime;",
				"import java.io.File;",
				"import javax.validation.constraints.Email;",
				"import javax.validation.constraints.NotNull;",
				"import javax.validation.constraints.NotBlank;",
				"import java.util.List;",
				"import java.util.Map;",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaJavaCodeGenerator{}
			sort.Strings(tt.want)
			if got := thiz.getImports(tt.args.messagePayload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMosaicKafkaJavaCodeGenerator_convertToJavaType(t *testing.T) {
	type args struct {
		property map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int32",
			args: args{
				property: map[string]interface{}{
					"type": "int32",
				},
			},
			want: "Integer",
		},
		{
			name: "int64",
			args: args{
				property: map[string]interface{}{
					"type": "int64",
				},
			},
			want: "Long",
		},
		{
			name: "string",
			args: args{
				property: map[string]interface{}{
					"type": "string",
				},
			},
			want: "String",
		},
		{
			name: "email",
			args: args{
				property: map[string]interface{}{
					"type":   "string",
					"format": "email",
				},
			},
			want: "String",
		},
		{
			name: "boolean",
			args: args{
				property: map[string]interface{}{
					"type": "boolean",
				},
			},
			want: "Boolean",
		},
		{
			name: "float",
			args: args{
				property: map[string]interface{}{
					"type": "float",
				},
			},
			want: "Float",
		},
		{
			name: "double",
			args: args{
				property: map[string]interface{}{
					"type": "double",
				},
			},
			want: "Double",
		},
		{
			name: "binary",
			args: args{
				property: map[string]interface{}{
					"type":   "string",
					"format": "binary",
				},
			},
			want: "File",
		},
		{
			name: "date",
			args: args{
				property: map[string]interface{}{
					"type": "date",
				},
			},
			want: "OffsetDateTime",
		},
		{
			name: "date-time",
			args: args{
				property: map[string]interface{}{
					"type": "date-time",
				},
			},
			want: "OffsetDateTime",
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
			want: "Map<String,String>",
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
			want: "List<TestObject>",
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
			want: "List<File>",
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
			want: "List<String>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaJavaCodeGenerator{}
			if got := thiz.convertToJavaType(tt.args.property); got != tt.want {
				t.Errorf("convertToJavaType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMosaicKafkaJavaCodeGenerator_getAnnotations(t *testing.T) {
	type args struct {
		propertyName string
		property     map[string]interface{}
		required     []interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "required default email",
			args: args{
				propertyName: "testProperty",
				property: map[string]interface{}{
					"default": "test@example.com",
					"type":    "string",
					"format":  "email",
				},
				required: []interface{}{
					"testProperty",
				},
			},
			want: []string{
				"@NotNull",
				"@Builder.Default",
				"@Email",
			},
		},
		{
			name: "required string",
			args: args{
				propertyName: "testProperty",
				property: map[string]interface{}{
					"type": "string",
				},
				required: []interface{}{
					"testProperty",
				},
			},
			want: []string{
				"@NotBlank",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thiz := &MosaicKafkaJavaCodeGenerator{}
			sort.Strings(tt.want)
			if got := thiz.getAnnotations(tt.args.propertyName, tt.args.property, tt.args.required); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}
