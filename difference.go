// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

import (
	"io"
	"strings"
)

// Slice is an alias for []any, representing a slice of any type.
type Slice = []any

// Map is an alias for map[string]any, representing a map with string keys and any values.
type Map = map[string]any

// JsonDifference represents the difference between two JSON structures.
type JsonDifference struct {
	diff *Map
}

// BetweenJson compares two JSON inputs and returns their difference.
//
// Parameters:
//   - expected: io.Reader containing the expected JSON data
//   - received: io.Reader containing the received JSON data
//
// Returns:
//   - *JsonDifference: A pointer to a JsonDifference struct containing the differences
//   - error: An error if there was a problem decoding the JSON inputs
//
// The function reads JSON data from both inputs, compares them, and generates a difference map.
func BetweenJson(expected, received io.Reader) (*JsonDifference, error) {
	diff := make(Map)
	expectedJson := make(Map)
	receivedJson := make(Map)

	if err := decodeJsonIntoMap(expected, &expectedJson); err != nil {
		return nil, err
	}

	if err := decodeJsonIntoMap(received, &receivedJson); err != nil {
		return nil, err
	}

	compareMaps(&diff, &expectedJson, &receivedJson, false)

	return &JsonDifference{&diff}, nil
}

// HasDifferences checks if there are any differences between the compared JSON structures.
//
// Returns:
//   - bool: true if differences exist, false otherwise
func (jd *JsonDifference) HasDifferences() bool {
	return len(*jd.diff) > 0
}

// GenerateReport creates a human-readable string report of the differences.
//
// Returns:
//   - string: A formatted string showing the differences between the expected and received JSON
//
// The report uses color coding and indentation to clearly display the differences:
//   - Red lines (prefixed with '-') indicate expected values
//   - Green lines (prefixed with '+') indicate received values
func (jd *JsonDifference) GenerateReport() string {
	var report strings.Builder

	report.WriteString("\n")
	report.WriteString(printRed("- Expected"))
	report.WriteString("\n")
	report.WriteString(printGreen("+ Received"))
	report.WriteString("\n")
	report.WriteString("\n")

	report.WriteString("  {")
	report.WriteString("\n")
	report.WriteString(printMapDiff(jd.diff, 1, nil))
	report.WriteString("\n")
	report.WriteString("  }")

	return report.String()
}
