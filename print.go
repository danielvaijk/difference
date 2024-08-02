package difference

import (
	"fmt"
	"strings"
)

func printMapDiff(diff *Map, indentationLevel int, sign string) string {
	var printBuilder strings.Builder

	propertyIndex := 0
	propertyPrints := make([]string, 0)

	for propertyKey, propertyValue := range *diff {
		suffix := ""
		isLastProperty := propertyIndex == len(*diff)-1

		if !isLastProperty {
			suffix = ","
		}

		propertyPrints = append(
			propertyPrints,
			printPropertyDiff(
				sign+propertyKey,
				propertyValue,
				indentationLevel,
				suffix,
			),
		)

		propertyIndex++
	}

	for propertyIndex, propertyPrint := range propertyPrints {
		printLines := strings.Split(propertyPrint, "\n")
		isLastProperty := propertyIndex == len(propertyPrints)-1

		for lineIndex, printLine := range printLines {
			isLastLine := lineIndex == len(printLines)-1

			if strings.HasPrefix(printLine, removed) {
				printBuilder.WriteString(printRed(printLine))
			} else if strings.HasPrefix(printLine, added) {
				printBuilder.WriteString(printGreen(printLine))
			} else {
				printBuilder.WriteString(printLine)
			}

			if !isLastProperty || !isLastLine {
				printBuilder.WriteString("\n")
			}
		}
	}

	return printBuilder.String()
}

func printPropertyDiff(
	key string,
	value any,
	indentationLevel int,
	suffix string,
) string {
	var printBuilder strings.Builder

	prefix, indentation := getPrefixSpacers(key, indentationLevel)

	sign := strings.TrimSpace(prefix)
	keyWithoutSign, _ := strings.CutPrefix(key, sign)
	propertyNames := strings.Split(keyWithoutSign, ".")
	propertyName := propertyNames[len(propertyNames)-1]

	printBuilder.WriteString(prefix)
	printBuilder.WriteString(indentation)
	printBuilder.WriteString(propertyName)
	printBuilder.WriteString(": ")

	switch value := (value).(type) {
	case Map:
		printBuilder.WriteString("{")
		printBuilder.WriteString("\n")
	case []Slice:
		printBuilder.WriteString("[")
		printBuilder.WriteString("\n")
	default:
		printBuilder.WriteString(formatValue(prefix+indentation, value))
		printBuilder.WriteString(suffix)
	}

	switch value := (value).(type) {
	case Map:
		printBuilder.WriteString(printMapDiff(&value, indentationLevel+1, sign))
	case []Slice:
		printBuilder.WriteString(printSliceDiff(&value, indentationLevel+1))
	}

	switch (value).(type) {
	case Map:
		printBuilder.WriteString("\n")
		printBuilder.WriteString(prefix)
		printBuilder.WriteString(indentation)
		printBuilder.WriteString("}")
		printBuilder.WriteString(suffix)
	case []Slice:
		printBuilder.WriteString("\n")
		printBuilder.WriteString(prefix)
		printBuilder.WriteString(indentation)
		printBuilder.WriteString("]")
		printBuilder.WriteString(suffix)
	}

	return printBuilder.String()
}

func printSliceDiff(slices *[]Slice, indentationLevel int) string {
	var printBuilder strings.Builder

	for index, pair := range *slices {
		if len(pair) != 2 {
			panic("malformed slice pair for slice diff")
		}

		prefix, indentation := getPrefixSpacers(
			pair[0].(string),
			indentationLevel,
		)

		printBuilder.WriteString(prefix)
		printBuilder.WriteString(indentation)
		printBuilder.WriteString(formatValue("", pair[1]))

		if index < len(*slices)-1 {
			printBuilder.WriteString(",")
			printBuilder.WriteString("\n")
		}
	}

	return printBuilder.String()
}

func printRed(str string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", str)
}

func printGreen(str string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", str)
}
