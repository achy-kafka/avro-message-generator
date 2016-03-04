package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
)

func GenerateMessage(schema string) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var schemaMap map[string]interface{}
	err := json.Unmarshal([]byte(schema), &schemaMap)
	must(err)

	parsedSchema := make(map[string]interface{})

	recursivelyPopulate(schemaMap, parsedSchema)

	parsedString, err := json.Marshal(parsedSchema)
	return string(parsedString)
}

func recursivelyPopulate(schema map[string]interface{}, parsed map[string]interface{}) {
	var valName string
	if schema["name"] != nil {
		valName = schema["name"].(string)
	}

	switch schema["type"].(type) {
	case []interface{}:
		handleUnion(schema, parsed, schema["type"].([]interface{}))
	case string:
		valType := schema["type"].(string)

		switch valType {
		case "record":
			fields := schema["fields"].([]interface{})
			parsed[valName] = make(map[string]interface{})
			for i := 0; i < len(fields); i++ {
				recursivelyPopulate(fields[i].(map[string]interface{}), parsed[valName].(map[string]interface{}))
			}
		case "array":
			subType := schema["items"]
			arrLen := rand.Intn(10)
			arr := make([]interface{}, arrLen)

			switch subType.(type) {
			case string:
				for i, _ := range arr {
					arr[i] = generateValueForType(subType.(string))
				}
			case []interface{}:
				for i, _ := range arr {
					arrMap := make(map[string]interface{})
					arr[i] = arrMap
					handleUnion(schema, arrMap, subType.([]interface{}))
				}
			default:
				for i, _ := range arr {
					arrMap := make(map[string]interface{})
					arr[i] = arrMap
					recursivelyPopulate(subType.(map[string]interface{}), arrMap)
				}
			}
			parsed[valName] = arr
		case "fixed":
			parsed[valName] = generateValueForType("bytes", int(schema["size"].(float64)))
		case "enum":
			symbols := schema["symbols"].([]interface{})
			index := rand.Intn(len(symbols))
			parsed[valName] = symbols[index].(string)
		case "map":
			switch schema["values"].(type) {
			case []interface{}:
				handleUnion(schema, parsed, schema["values"].([]interface{}))
			default:
				mapType := schema["values"].(string)
				subMap := make(map[string]interface{})
				numKeys := rand.Intn(10)
				for i := 0; i < numKeys; i++ {
					key := generateValueForType("string").(string)
					subMap[key] = generateValueForType(mapType)
				}
				parsed[valName] = subMap
			}

		default:
			parsed[valName] = generateValueForType(valType)
		}
	default:
		valType := schema["type"].(map[string]interface{})
		parsed[valName] = make(map[string]interface{})
		recursivelyPopulate(valType, parsed[valName].(map[string]interface{}))
	}

}

func handleUnion(schema map[string]interface{}, parsed map[string]interface{}, valType []interface{}) {
	var valName string
	if schema["name"] != nil {
		valName = schema["name"].(string)
	}
	index := rand.Intn(len(valType))
	selectedType := valType[index]
	switch reflect.TypeOf(selectedType).String() {
	case "string":
		subMap := make(map[string]interface{})
		stringType := selectedType.(string)
		subMap[stringType] = generateValueForType(stringType)
		parsed[valName] = subMap
	default:
		parsed[valName] = make(map[string]interface{})
		subMap := make(map[string]interface{})
		typeMap := selectedType.(map[string]interface{})
		typeKey := typeMap["type"].(string)
		subMap[typeKey] = make(map[string]interface{})
		parsed[valName] = subMap
		recursivelyPopulate(selectedType.(map[string]interface{}), subMap[typeKey].(map[string]interface{}))
	}

}

func generateValueForType(valueType string, length ...int) interface{} {
	genLength := 10
	if len(length) > 0 {
		genLength = length[0]
	}
	switch valueType {
	case "string":
		p := make([]byte, genLength)
		rand.Read(p)
		return base64.URLEncoding.EncodeToString(p)
	case "int":
		return rand.Int31()
	case "long":
		return rand.Int63()
	case "float":
		return rand.Float32()
	case "double":
		return rand.Float64()
	case "bytes":
		p := make([]byte, genLength)
		rand.Read(p)
		return p
	case "boolean":
		return rand.Intn(2) != 0
	}
	return nil
}
