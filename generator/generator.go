package generator

import (
	"github.com/asyncapi/parser-go/pkg/jsonpath"
	v2 "github.com/asyncapi/parser-go/pkg/parser/v2"
)

type Generator interface {
	Generate(asyncApiSpecPath string, out string) (string, error)
}

func LoadAsyncApiSpecWithParser(asyncApiSpecPath string) (map[string]interface{}, error) {
	refloader := jsonpath.NewRefLoader(nil)
	parser := v2.NewParser(refloader)
	data, err := parser.Load(asyncApiSpecPath)
	if err != nil {
		return nil, err
	}
	err = parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetMessages(data map[string]interface{}) []map[string]interface{} {
	var out []map[string]interface{}
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["publish"] != nil {
			operationProperties := operation["publish"].(map[string]interface{})
			message := operationProperties["message"].(map[string]interface{})
			out = append(out, message)
		} else if operation["subscribe"] != nil {
			operationProperties := operation["subscribe"].(map[string]interface{})
			message := operationProperties["message"].(map[string]interface{})
			out = append(out, message)
		}
	}
	return out
}

func CheckRequired(propertyName string, required []interface{}) bool {
	for _, requiredProperty := range required {
		reqProp := requiredProperty.(string)
		if propertyName == reqProp {
			return true
		}
	}
	return false
}

func GetNestedObjects(messages []map[string]interface{}) []map[string]interface{} {
	var nestedObject []map[string]interface{}
	for _, message := range messages {
		pay := message["payload"].(map[string]interface{})
		for _, prop := range pay["properties"].(map[string]interface{}) {
			property := prop.(map[string]interface{})
			if property["type"] == "object" {
				if property["additionalProperties"] != nil {
					apb, ok := property["additionalProperties"].(bool)
					if ok {
						if !apb {
							property["parent"] = message["name"]
							nestedObject = append(nestedObject, property)
						}
					}
				} else {
					property["parent"] = message["name"]
					nestedObject = append(nestedObject, property)
				}

			}
		}
	}
	return nestedObject
}

func GetItemObjects(messages []map[string]interface{}) []map[string]interface{} {
	var nestedObject []map[string]interface{}
	for _, message := range messages {
		pay := message["payload"].(map[string]interface{})
		for _, prop := range pay["properties"].(map[string]interface{}) {
			property := prop.(map[string]interface{})
			if property["type"] == "array" {
				if property["additionalProperties"] != nil {
					apb, ok := property["additionalProperties"].(bool)
					if ok {
						if !apb {
							propItems := property["items"].(map[string]interface{})
							if propItems["type"] == "object" {
								property["parent"] = message["name"]
								nestedObject = append(nestedObject, property)
							}
						}
					}
				} else {
					propItems := property["items"].(map[string]interface{})
					if propItems["type"] == "object" {
						property["parent"] = message["name"]
						nestedObject = append(nestedObject, property)
					}
				}

			}
		}
	}
	return nestedObject
}

func HasDefault(property map[string]interface{}) bool {
	return property["default"] != nil
}

func HasProducer(data map[string]interface{}) bool {
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["subscribe"] != nil {
			return true
		}
	}
	return false
}

func HasConsumer(data map[string]interface{}) bool {
	for _, channel := range data["channels"].(map[string]interface{}) {
		operation := channel.(map[string]interface{})
		if operation["publish"] != nil {
			return true
		}
	}
	return false
}
