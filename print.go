package difference

import (
	"fmt"
	"strings"
)

func printMapDiff(diff Map, indentationLevel int) string {
	var propertyPrints []string
	var printBuilder strings.Builder

	propertyIndex := 0

	for propertyKey, propertyValue := range diff {
		isLastProperty := propertyIndex == len(diff)-1
		prefix, tab := getPrefixSpacers(propertyKey, indentationLevel)
		isAdditionOrRemoval := len(prefix) > 0

		if isAdditionOrRemoval {
			propertyPrints = append(
				propertyPrints,
				printWholePropertyDiff(
					propertyKey,
					propertyValue,
					isLastProperty,
					prefix,
					tab,
				),
			)
		} else {
			propertyPrints = append(
				propertyPrints,
				printPartialPropertyDiff(
					propertyKey,
					propertyValue,
					indentationLevel,
				),
			)
		}

		propertyIndex++
	}

	for propertyIndex, propertyPrint := range propertyPrints {
		printLines := strings.Split(propertyPrint, "\n")
		isLastProperty := propertyIndex == len(propertyPrints)-1

		for lineIndex, printLine := range printLines {
			suffix := ""
			isLastLine := lineIndex == len(printLines)-1

			if !isLastProperty && isLastLine {
				suffix = ","
			}

			if strings.HasPrefix(printLine, "-") {
				printBuilder.WriteString(printRed(printLine))
			} else if strings.HasPrefix(printLine, "+") {
				printBuilder.WriteString(printGreen(printLine))
			} else {
				printBuilder.WriteString(printLine + suffix)
			}

			if !isLastProperty || !isLastLine {
				printBuilder.WriteString("\n")
			}
		}
	}

	return printBuilder.String()
}

func printWholePropertyDiff(
	key string,
	value any,
	isLast bool,
	prefix string,
	tab string,
) string {
	var printBuilder strings.Builder

	for index, line := range strings.Split(formatValue(value), "\n") {
		printBuilder.WriteString(prefix + tab)

		propertyNames := strings.Split(key[1:], ".")
		propertyName := propertyNames[len(propertyNames)-1]

		suffix := ""

		if !isLast {
			suffix = ","
		}

		if index == 0 {
			printBuilder.WriteString(propertyName + ": " + line + suffix)
		} else {
			printBuilder.WriteString(line + suffix)
		}
	}

	return printBuilder.String()
}

func printPartialPropertyDiff(
	key string,
	value any,
	indentationLevel int,
) string {
	var openBracket string
	var closeBracket string

	var printBuilder strings.Builder

	switch (value).(type) {
	case Map:
		openBracket = "{"
		closeBracket = "}"
	case []Slice:
		openBracket = "["
		closeBracket = "]"
	}

	_, tab := getPrefixSpacers(key, indentationLevel)

	printBuilder.WriteString("   " + tab + key + ": " + openBracket + "\n")

	switch value := (value).(type) {
	case Map:
		printBuilder.WriteString(printMapDiff(value, indentationLevel+1))
	case []Slice:
		printBuilder.WriteString(printSliceDiff(value, indentationLevel+1))
	}

	printBuilder.WriteString("\n   " + tab + closeBracket)

	return printBuilder.String()
}

func printSliceDiff(slices []Slice, indentationLevel int) string {
	var printBuilder strings.Builder

	for index, pair := range slices {
		if len(pair) != 2 {
			panic("malformed slice pair for slice diff")
		}

		suffix := ""

		if index < len(slices)-1 {
			suffix = ",\n"
		}

		prefix, tab := getPrefixSpacers(pair[0].(string), indentationLevel)

		if len(prefix) > 0 {
			printBuilder.WriteString(prefix)
		} else {
			printBuilder.WriteString("   ")
		}

		printBuilder.WriteString(tab + formatValue(pair[1]) + suffix)
	}

	return printBuilder.String()
}

func printRed(str string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", str)
}

func printGreen(str string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", str)
}
