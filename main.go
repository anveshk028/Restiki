package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var nilValue *interface{} = nil  // Explicitly using a nil interface pointer
var trueValue interface{} = true // Explicitly using a true boolean interface

func main() {
	// Specify the path to your JSON file
	filePath := "input.json" // Adjust the file path as necessary

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the contents of the file
	inputData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Initialize
	outputData := make(map[string]*interface{})
	var jsonInputData map[string]interface{}
	err = json.Unmarshal([]byte(inputData), &jsonInputData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Handle each schemaless key with S, N, BOOL, NULL, M, L struct mapping
	for key, value := range jsonInputData {

		if numString, err := IsNumberInterface(value); err == nil && numString != "" {
			fmt.Printf("\nkey = %s \n numberInterface = %v", key, numString)
			parsedFloat64, err := ParseFloat64(numString)
			if err == nil {
				parsedFloat64Interface := interface{}(parsedFloat64)
				outputData[key] = &parsedFloat64Interface
			}
		} else if str, err := IsStringInterface(value); err == nil && str != "" {
			fmt.Printf("\nkey = %s \n stringInterface = %v", key, str)
			parsedDate, err := ParseDateRFC3339ToUnix(str)
			if err == nil {
				parsedDateInterface := interface{}(parsedDate)
				outputData[key] = &parsedDateInterface
			} else {
				strInterface := interface{}(str)
				outputData[key] = &strInterface
			}
		} else if boolString, err := IsBooleanInterface(value); err == nil && boolString != "" {
			fmt.Printf("\nkey = %s \n booleanInterface = %v", key, boolString)
			if ParseBoolean(boolString) != nil {
				outputData[key] = &trueValue
			}
		} else if nullString, err := IsNullInterface(value); err == nil && nullString != "" {
			fmt.Printf("\nkey = %s \n nullInterface = %v", key, nullString)
			if ParseBoolean(nullString) != nil {
				outputData[key] = nilValue
			}
		} else {
			// Solve for M & L
		}
	}

	jsonOutputData, err := json.MarshalIndent(outputData, " ", "  ") // indent with four spaces
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println("\nJSON Output:")
	fmt.Println(string(jsonOutputData))
}

type NumberInterface struct {
	N string `json:"N,omitempty"`
}

type StringInterface struct {
	S string `json:"S,omitempty"`
}

type BooleanInterface struct {
	BOOL string `json:"BOOL,omitempty"`
}

type NullInterface struct {
	NULL string `json:"NULL,omitempty"`
}

type MapInterface struct {
	M map[string]interface{} `json:"M,omitempty"`
}

type ListInterface struct {
	L []interface{} `json:"L,omitempty"`
}

func IsNumberInterface(data interface{}) (string, error) {
	var numInt NumberInterface
	temp, err := json.Marshal(data)
	err = json.Unmarshal(temp, &numInt)
	return TrimSpace(numInt.N), err
}

func IsStringInterface(data interface{}) (string, error) {
	var strInt StringInterface
	temp, err := json.Marshal(data)
	err = json.Unmarshal(temp, &strInt)
	return TrimSpace(strInt.S), err
}

func IsBooleanInterface(data interface{}) (string, error) {
	var boolInt BooleanInterface
	temp, err := json.Marshal(data)
	err = json.Unmarshal(temp, &boolInt)
	return TrimSpace(boolInt.BOOL), err
}

func IsNullInterface(data interface{}) (string, error) {
	var nullInt NullInterface
	temp, err := json.Marshal(data)
	err = json.Unmarshal(temp, &nullInt)
	return TrimSpace(nullInt.NULL), err
}

func IsMapInterface(data interface{}) (*MapInterface, error) {
	var mapInt MapInterface
	temp, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temp, &mapInt)
	fmt.Printf("data = %v mapInteface = %v", data, mapInt)
	return &mapInt, err
}

func IsListInterface(data interface{}) (*ListInterface, error) {
	var listInt ListInterface
	temp, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temp, &listInt)
	return &listInt, err
}

func TrimSpace(value string) string {
	return strings.TrimSpace(value)
}

func ParseFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func ParseDateRFC3339ToUnix(value string) (int64, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0, err
	}
	return t.Unix(), err
}

func ParseBoolean(value string) *string {
	// Switch on the value to determine the boolean output
	switch value {
	case "1", "t", "T", "TRUE", "true", "True":
		return &value
	case "0", "f", "F", "FALSE", "false", "False":
		return nil
	default:
		return nil
		// **Must** omit fields with invalid `Boolean` values
	}
}
