// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

import (
	"fmt"
	"strings"
)

const tab = "  "

func formatValue(prefix string, value any) string {
	switch value := value.(type) {
	case Map:
		return formatMap(prefix, value)
	case Slice:
		return formatSlice(prefix, value)
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("%q", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatMap(prefix string, value Map) string {
	var index int
	var output strings.Builder

	output.WriteString("{")
	output.WriteString("\n")

	for propertyName, propertyValue := range value {
		output.WriteString(prefix)
		output.WriteString(tab)
		output.WriteString(propertyName)
		output.WriteString(": ")
		output.WriteString(formatValue(prefix, propertyValue))

		if index < len(value)-1 {
			output.WriteString(",")
			output.WriteString("\n")
		}

		index++
	}

	output.WriteString("\n")
	output.WriteString(prefix)
	output.WriteString("}")

	return output.String()
}

func formatSlice(prefix string, slice []any) string {
	var output strings.Builder

	output.WriteString("[")
	output.WriteString("\n")

	for index, value := range slice {
		output.WriteString(prefix)
		output.WriteString(tab)
		output.WriteString(formatValue(prefix+tab, value))

		if index < len(slice)-1 {
			output.WriteString(",")
			output.WriteString("\n")
		}
	}

	output.WriteString("\n")
	output.WriteString(prefix)
	output.WriteString("]")

	return output.String()
}

func getPrefixSpacers(sign rune, indentationLevel int) (string, string) {
	return string(sign) + tab, strings.Repeat(tab, indentationLevel)
}
