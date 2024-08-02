package indifferent

import (
	"io"
	"reflect"
)

const (
	added   = "+"
	removed = "-"
	common  = " "
)

func CompareJson(expected, received io.Reader) (Map, error) {
	var expectedJson Map
	var receivedJson Map

	if err := decodeJsonIntoMap(expected, &expectedJson); err != nil {
		return nil, err
	}

	if err := decodeJsonIntoMap(received, &receivedJson); err != nil {
		return nil, err
	}

	diff := make(Map)
	compareMaps(&diff, &expectedJson, &receivedJson)

	return diff, nil
}

func compareMaps(diff, expected, received *Map, propertyPath ...string) {
	if len(propertyPath) == 0 {
		propertyPath = append(propertyPath, "")
	}

	for expectedKey, expectedValue := range *expected {
		receivedValue, wasFound := (*received)[expectedKey]
		expectedKey = propertyPath[0] + expectedKey

		if !wasFound {
			registerPropertyRemoval(diff, expectedKey, expectedValue)
			continue
		}

		expectedType := reflect.TypeOf(expectedValue)
		receivedType := reflect.TypeOf(receivedValue)

		if expectedType != receivedType {
			registerPropertyChange(diff, expectedKey, expectedValue, receivedValue)
			continue
		}

		switch expectedValue := expectedValue.(type) {
		case Map:
			nestedDiff := make(Map)

			compareMaps(
				&nestedDiff,
				&expectedValue,
				receivedValue.(*Map),
				expectedKey+".",
			)

			if len(nestedDiff) > 0 {
				registerNestedDiff(diff, expectedKey, &nestedDiff)
			}
		case Slice:
			if !reflect.DeepEqual(expectedValue, receivedValue) {
				registerPropertyChange(diff, expectedKey, expectedValue, receivedValue)
			}
		default:
			if expectedValue != receivedValue {
				registerPropertyChange(diff, expectedKey, expectedValue, receivedValue)
			}
		}
	}

	for receivedKey, receivedValue := range *received {
		if _, wasFound := (*expected)[receivedKey]; !wasFound {
			registerPropertyAddition(diff, propertyPath[0]+receivedKey, receivedValue)
		}
	}
}

func registerPropertyRemoval(diff *Map, propertyKey string, propertyValue any) {
	(*diff)[removed+propertyKey] = propertyValue
}

func registerPropertyAddition(diff *Map, propertyKey string, propertyValue any) {
	(*diff)[added+propertyKey] = propertyValue
}

func registerPropertyChange(diff *Map, propertyKey string, expectedValue, receivedValue any) {
	(*diff)[removed+propertyKey] = expectedValue
	(*diff)[added+propertyKey] = receivedValue
}

func registerNestedDiff(diff *Map, propertyKey string, nestedDiff *Map) {
	(*diff)[common+propertyKey] = nestedDiff
}
