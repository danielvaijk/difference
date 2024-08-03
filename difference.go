package difference

import (
	"io"
	"strings"
)

type Slice = []any
type Map = map[string]any

func BetweenJson(expected, received io.Reader) (Map, error) {
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

func GenerateReport(diff Map) string {
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
