package difference

import (
	"reflect"
)

func compareMaps(diff, expected, received *Map, propertyPath ...string) {
	if len(propertyPath) == 0 {
		propertyPath = append(propertyPath, "")
	}

	for expectedKey, expectedValue := range *expected {
		receivedValue, wasFound := (*received)[expectedKey]
		expectedKey = propertyPath[0] + expectedKey

		if !wasFound {
			registerRemovedProperty(diff, expectedKey, expectedValue)
			continue
		}

		expectedType := reflect.TypeOf(expectedValue)
		receivedType := reflect.TypeOf(receivedValue)

		if expectedType != receivedType {
			registerChangedProperty(diff, expectedKey, expectedValue, receivedValue)
			continue
		}

		switch expectedValue := expectedValue.(type) {
		case Map:
			mapDiff := make(Map)
			receivedMap := receivedValue.(Map)

			compareMaps(
				&mapDiff,
				&expectedValue,
				&receivedMap,
				expectedKey+".",
			)

			if len(mapDiff) > 0 {
				registerNestedDiff(diff, expectedKey, mapDiff)
			}
		case Slice:
			sliceDiff := make([]Slice, 0)
			receivedSlice := receivedValue.(Slice)

			compareSlices(
				&sliceDiff,
				&expectedValue,
				&receivedSlice,
			)

			if len(sliceDiff) > 0 {
				registerNestedDiff(diff, expectedKey, sliceDiff)
			}
		default:
			if expectedValue != receivedValue {
				registerChangedProperty(diff, expectedKey, expectedValue, receivedValue)
			}
		}
	}

	for receivedKey, receivedValue := range *received {
		if _, wasFound := (*expected)[receivedKey]; !wasFound {
			registerAddedProperty(diff, propertyPath[0]+receivedKey, receivedValue)
		}
	}
}

func compareSlices(diff *[]Slice, expected, received *Slice) {
	expectedIndex := 0
	receivedIndex := 0

	for {
		isExpectedOutOfBounds := expectedIndex > len(*expected)-1
		isReceivedOutOfBounds := receivedIndex > len(*received)-1

		if isReceivedOutOfBounds {
			break
		}

		if !isExpectedOutOfBounds && (*expected)[expectedIndex] == (*received)[receivedIndex] {
			registerMutualValue(diff, (*expected)[expectedIndex])
			expectedIndex++
			receivedIndex++
			continue
		}

		// Any differences at the start are considered removals.
		if receivedIndex == 0 {
			registerRemovedValue(diff, (*expected)[expectedIndex])
			expectedIndex++
			continue
		}

		// Any differences at the end are considered additions.
		if isExpectedOutOfBounds {
			registerAddedValue(diff, (*received)[receivedIndex])
			receivedIndex++
			continue
		}

		// Any differences between the start and end are considered changes.
		registerChangedValue(diff, (*expected)[expectedIndex], (*received)[receivedIndex])

		expectedIndex++
		receivedIndex++
	}
}

func registerRemovedProperty(diff *Map, propertyKey string, propertyValue any) {
	(*diff)[string(removed)+propertyKey] = propertyValue
}

func registerAddedProperty(diff *Map, propertyKey string, propertyValue any) {
	(*diff)[string(added)+propertyKey] = propertyValue
}

func registerChangedProperty(diff *Map, propertyKey string, expectedValue, receivedValue any) {
	registerRemovedProperty(diff, propertyKey, expectedValue)
	registerAddedProperty(diff, propertyKey, receivedValue)
}

func registerNestedDiff(diff *Map, propertyKey string, nestedDiff any) {
	(*diff)[string(unchanged)+propertyKey] = nestedDiff
}

func registerRemovedValue(diff *[]Slice, value any) {
	*diff = append(*diff, Slice{removed, value})
}

func registerAddedValue(diff *[]Slice, value any) {
	*diff = append(*diff, Slice{added, value})
}

func registerChangedValue(diff *[]Slice, expectedValue, receivedValue any) {
	registerRemovedValue(diff, expectedValue)
	registerAddedValue(diff, receivedValue)
}

func registerMutualValue(diff *[]Slice, mutualValue any) {
	*diff = append(*diff, Slice{unchanged, mutualValue})
}
