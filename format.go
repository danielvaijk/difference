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

func getPrefixSpacers(key string, indentationLevel int) (string, string) {
	indentation := strings.Repeat(tab, indentationLevel)

	if strings.HasPrefix(key, removed) {
		return removed + "  ", indentation
	} else if strings.HasPrefix(key, added) {
		return added + "  ", indentation
	} else {
		return "   ", indentation
	}
}
