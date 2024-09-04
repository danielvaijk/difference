// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

import (
	"io"
	"strings"
)

type Slice = []any
type Map = map[string]any

type JsonDifference struct {
	diff *Map
}

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

func (jd *JsonDifference) HasDifferences() bool {
	return len(*jd.diff) > 0
}

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
