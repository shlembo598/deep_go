package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var result []string

	vData := reflect.ValueOf(person)
	tData := vData.Type()

	if vData.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < vData.NumField(); i++ {
		fieldType := tData.Field(i)
		tag := fieldType.Tag.Get("properties")
		if tag == "" {
			continue
		}

		tagName, hasOmitEmpty := parseTag(tag)
		if tagName == "" {
			continue
		}

		field := vData.Field(i)
		isZeroValue := checkValueIsZero(field)

		if isZeroValue && hasOmitEmpty {
			continue
		}

		tagValue := serializeValue(field)
		result = append(result, fmt.Sprintf("%s=%v", tagName, tagValue))
	}

	return strings.Join(result, "\n")
}

func parseTag(tag string) (tagName string, hasOmitEmpty bool) {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "omitempty" {
			hasOmitEmpty = true
		} else if tagName == "" {
			tagName = part
		}
	}
	return
}

func checkValueIsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int:
		return v.Int() == 0
	case reflect.Bool:
		return !v.Bool()
	default:
		return v.IsZero()
	}
}

func serializeValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	default:
		return ""
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
