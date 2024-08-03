package difference

import (
	"fmt"
	"strings"
)

func printMapDiff(diff *Map, indentationLevel int, sign string) string {
	var printBuilder strings.Builder

	diffIndex := 0
	diffCount := len(*diff)

	prints := make([]string, 0)
	previousPrint := ""

	for propertyKey, propertyValue := range *diff {
		shouldAddComma := true

		isLastPrint := diffIndex == diffCount-1
		hasPreviousPrint := len(previousPrint) > 0

		currentPrint := printPropertyDiff(
			sign+propertyKey,
			propertyValue,
			indentationLevel,
		)

		if hasPreviousPrint && isLastPrint {
			currentSign := currentPrint[:1]
			previousSign := previousPrint[:1]

			shouldAddComma = !isChangeSet(
				previousSign,
				currentSign,
			)
		}

		if hasPreviousPrint {
			if shouldAddComma {
				prints = append(prints, previousPrint+",")
			} else {
				prints = append(prints, previousPrint)
			}
		}

		previousPrint = currentPrint
		diffIndex++
	}

	prints = append(prints, previousPrint)

	for printIndex, printResult := range prints {
		printLines := strings.Split(printResult, "\n")
		isLastPrint := printIndex == len(prints)-1

		for lineIndex, printLine := range printLines {
			isLastLine := lineIndex == len(printLines)-1

			if strings.HasPrefix(printLine, removed) {
				printBuilder.WriteString(printRed(printLine))
			} else if strings.HasPrefix(printLine, added) {
				printBuilder.WriteString(printGreen(printLine))
			} else {
				printBuilder.WriteString(printLine)
			}

			if !isLastPrint || !isLastLine {
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
	case []Slice:
		printBuilder.WriteString("\n")
		printBuilder.WriteString(prefix)
		printBuilder.WriteString(indentation)
		printBuilder.WriteString("]")
	}

	return printBuilder.String()
}

func printSliceDiff(slices *[]Slice, indentationLevel int) string {
	var printBuilder strings.Builder

	totalCount := len(*slices)

	for index, pair := range *slices {
		if len(pair) != 2 {
			panic("malformed slice pair for slice diff")
		}

		currentSign := pair[0].(string)
		prefix, indentation := getPrefixSpacers(currentSign, indentationLevel)

		printBuilder.WriteString(prefix)
		printBuilder.WriteString(indentation)
		printBuilder.WriteString(formatValue("", pair[1]))

		shouldAddComma := true

		isLast := index == totalCount-1
		isBeforeLast := index == totalCount-2

		if isLast {
			shouldAddComma = false
		} else if isBeforeLast {
			upcomingSign := (*slices)[index+1][0].(string)
			shouldAddComma = !isChangeSet(currentSign, upcomingSign)
		}

		if shouldAddComma {
			printBuilder.WriteString(",")
		}

		if !isLast {
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
