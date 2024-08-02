package difference

import (
	"fmt"
	"strings"
)

func formatValue(value any) string {
	switch value := value.(type) {
	case string:
		return fmt.Sprintf("%q", value)
	case Map:
		return formatMap(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatMap(m Map) string {
	var pairs []string

	for fieldName, fieldValue := range m {
		pairs = append(pairs, "  "+fieldName+": "+formatValue(fieldValue))
	}

	return "{\n" + strings.Join(pairs, ",\n") + "\n}"
}

func getPrefixSpacers(key string, indentationLevel int) (string, string) {
	tab := strings.Repeat("  ", indentationLevel)

	if strings.HasPrefix(key, removed) {
		return removed + "  ", tab
	} else if strings.HasPrefix(key, added) {
		return added + "  ", tab
	} else {
		return "", tab
	}
}
