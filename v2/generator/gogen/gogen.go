package gogen

import (
	"embed"
	"fmt"
	"golang.org/x/exp/slices"
)

var typeConversionGoMap = map[string]string{
	"integer":   "int",
	"int32":     "int32",
	"int64":     "int64",
	"string":    "string",
	"email":     "string",
	"boolean":   "bool",
	"number":    "float32",
	"float":     "float32",
	"double":    "float64",
	"binary":    "[]byte",
	"date":      "time.Time",
	"date-time": "time.Time",
}

//go:embed templates
var templateFiles embed.FS

func getImports(data map[string]interface{}) []string {
	var out []string
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["publish"] != nil {
			operationProperties := operation["publish"].(map[string]interface{})
			imp := extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		} else if operation["subscribe"] != nil {
			operationProperties := operation["subscribe"].(map[string]interface{})
			imp := extractImport(operationProperties)
			if imp != "" {
				if !slices.Contains(out, imp) {
					out = append(out, imp)
				}
			}
		}

	}
	return out
}

func extractImport(operationProperties map[string]interface{}) string {
	message := operationProperties["message"].(map[string]interface{})
	payload := message["payload"].(map[string]interface{})
	properties := payload["properties"].(map[string]interface{})
	for _, property := range properties {
		prop := property.(map[string]interface{})
		if prop["format"] != nil {
			frm := prop["format"].(string)
			switch frm {
			case "date-time", "date":
				return "time"
			}
		}
	}

	return ""
}

func convertToGoType(property map[string]interface{}) string {
	switch property["type"] {
	case "object":
		if property["additionalProperties"] != nil {
			additionalProperties, ok := property["additionalProperties"].(map[string]interface{})
			if ok {
				if additionalProperties["type"] == "string" {
					return "map[string]string"
				}
			}
			_, ok = property["additionalProperties"].(bool)
			if ok {
				return property["title"].(string)
			}
		} else {
			return property["title"].(string)
		}
	case "array":
		typ := "[]"
		if property["items"] != nil {
			items := property["items"].(map[string]interface{})
			if items["format"] != nil {
				typ = typ + typeConversionGoMap[items["format"].(string)]
			} else if items["type"] == "object" {
				typ = typ + items["title"].(string)
			} else {
				typ = typ + typeConversionGoMap[items["type"].(string)]
			}
			return typ
		}
	default:
		if property["format"] != nil {
			return typeConversionGoMap[property["format"].(string)]
		} else {
			return typeConversionGoMap[property["type"].(string)]
		}
	}
	return ""
}

func validations(property map[string]interface{}, required bool) string {
	out := ` validate:"`
	if required && property["type"] != "boolean" {
		out = out + "required,"
	}
	if property["format"] != nil {
		switch property["format"].(string) {
		case "email":
			out = out + "email,"
		}
	}
	if property["minimum"] != nil {
		min := property["minimum"].(int)
		out = out + fmt.Sprintf("gte=%d,", min)
	} else if property["exclusiveMinimum"] != nil {
		min := property["exclusiveMinimum"].(int)
		out = out + fmt.Sprintf("gt=%d,", min)
	}
	if property["maximum"] != nil {
		max := property["maximum"].(int)
		out = out + fmt.Sprintf("lte=%d,", max)
	} else if property["exclusiveMaximum"] != nil {
		max := property["exclusiveMaximum"].(int)
		out = out + fmt.Sprintf("lt=%d,", max)
	}
	if property["type"] == "string" {
		if property["minLength"] != nil {
			min := property["minLength"].(int)
			out = out + fmt.Sprintf("min=%d,", min)
		}
		if property["maxLength"] != nil {
			max := property["maxLength"].(int)
			out = out + fmt.Sprintf("max=%d,", max)
		}
	}
	if property["type"] == "array" {
		if property["minItems"] != nil {
			min := property["minItems"].(int)
			out = out + fmt.Sprintf("min=%d,", min)
		}
		if property["maxItems"] != nil {
			max := property["maxItems"].(int)
			out = out + fmt.Sprintf("max=%d,", max)
		}
		if property["uniqueItems"] != nil {
			if property["uniqueItems"].(bool) {
				out = out + "unique,"
			}
		}
	}
	out = out[:len(out)-1] + `"`
	if out != ` validate:"` {
		return out
	}
	return ""
}
