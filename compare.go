// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

import (
	"reflect"
	"slices"
)

func compareMaps(diff, expected, received *Map, includeCommon bool, propertyPath ...string) bool {
	hasDifferences := false

	if len(propertyPath) == 0 {
		propertyPath = append(propertyPath, "")
	}

	for expectedKey, expectedValue := range *expected {
		receivedValue, wasFound := (*received)[expectedKey]
		expectedKey = propertyPath[0] + expectedKey

		if !wasFound {
			hasDifferences = true
			registerRemovedProperty(diff, expectedKey, expectedValue)
			continue
		}

		expectedType := reflect.TypeOf(expectedValue)
		receivedType := reflect.TypeOf(receivedValue)

		if expectedType != receivedType {
			hasDifferences = true
			registerChangedProperty(diff, expectedKey, expectedValue, receivedValue)
			continue
		}

		switch expectedValue := expectedValue.(type) {
		case Map:
			mapDiff := make(Map)
			receivedMap := receivedValue.(Map)

			areMapsDifferent := compareMaps(
				&mapDiff,
				&expectedValue,
				&receivedMap,
				includeCommon,
				expectedKey+".",
			)

			if areMapsDifferent {
				hasDifferences = true
				registerNestedPropertyDiff(diff, expectedKey, mapDiff)
			}
		case Slice:
			sliceDiff := make([]Slice, 0)
			receivedSlice := receivedValue.(Slice)

			areSlicesDifferent := compareSlices(
				&sliceDiff,
				&expectedValue,
				&receivedSlice,
			)

			if areSlicesDifferent {
				hasDifferences = true
				registerNestedPropertyDiff(diff, expectedKey, sliceDiff)
			}
		default:
			if expectedValue != receivedValue {
				hasDifferences = true
				registerChangedProperty(diff, expectedKey, expectedValue, receivedValue)
			} else if includeCommon {
				registerMutualProperty(diff, expectedKey, expectedValue)
			}
		}
	}

	for receivedKey, receivedValue := range *received {
		if _, wasFound := (*expected)[receivedKey]; !wasFound {
			hasDifferences = true
			registerAddedProperty(diff, propertyPath[0]+receivedKey, receivedValue)
		}
	}

	return hasDifferences
}

func compareSlices(diff *[]Slice, expectedPtr, receivedPtr *Slice) bool {
	expectedIndex := 0
	receivedIndex := 0

	expected := *expectedPtr
	received := *receivedPtr

	hasDifferences := false

	for expectedIndex < len(expected) || receivedIndex < len(received) {
		if expectedIndex < len(expected) && receivedIndex < len(received) {
			switch expectedValue := expected[expectedIndex].(type) {
			case Map:
				mapDiff := make(Map)
				receivedMap := received[receivedIndex].(Map)

				areMapsDifferent := compareMaps(
					&mapDiff,
					&expectedValue,
					&receivedMap,
					true,
				)

				if areMapsDifferent {
					hasDifferences = true
					registerNestedDiffValue(diff, mapDiff)
				} else {
					registerMutualValue(diff, expectedValue)
				}

				expectedIndex++
				receivedIndex++
			case Slice:
				sliceDiff := make([]Slice, 0)
				receivedSlice := received[receivedIndex].(Slice)

				areSlicesDifferent := compareSlices(
					&sliceDiff,
					&expectedValue,
					&receivedSlice,
				)

				if areSlicesDifferent {
					hasDifferences = true
					registerNestedDiffValue(diff, sliceDiff)
				} else {
					registerMutualValue(diff, expectedValue)
				}

				expectedIndex++
				receivedIndex++
			default:
				if expectedValue == received[receivedIndex] {
					// Both values are the same.
					registerMutualValue(diff, expectedValue)
					expectedIndex++
					receivedIndex++
				} else if slices.Contains(received[receivedIndex:], expected[expectedIndex]) {
					// The current expected value is received later: show the current received value as superfluous.
					registerAddedValue(diff, received[receivedIndex])
					hasDifferences = true
					receivedIndex++
				} else {
					// The current expected value is missing.
					registerRemovedValue(diff, expected[expectedIndex])
					hasDifferences = true
					expectedIndex++
				}
			}
		} else if expectedIndex < len(expected) {
			// We ran out of received values but expected more: show them as missing.
			registerRemovedValue(diff, expected[expectedIndex])
			hasDifferences = true
			expectedIndex++
		} else if receivedIndex < len(received) {
			// We ran out of expected values but received more: show them as superfluous.
			registerAddedValue(diff, received[receivedIndex])
			hasDifferences = true
			receivedIndex++
		}
	}

	return hasDifferences
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

func registerMutualProperty(diff *Map, propertyKey string, propertyValue any) {
	(*diff)[string(unchanged)+propertyKey] = propertyValue
}

func registerNestedPropertyDiff(diff *Map, propertyKey string, nestedDiff any) {
	(*diff)[string(unchanged)+propertyKey] = nestedDiff
}

func registerRemovedValue(diff *[]Slice, value any) {
	*diff = append(*diff, Slice{removed, value})
}

func registerAddedValue(diff *[]Slice, value any) {
	*diff = append(*diff, Slice{added, value})
}

func registerMutualValue(diff *[]Slice, mutualValue any) {
	*diff = append(*diff, Slice{unchanged, mutualValue})
}

func registerNestedDiffValue(diff *[]Slice, valueDiff any) {
	*diff = append(*diff, Slice{nested, valueDiff})
}
