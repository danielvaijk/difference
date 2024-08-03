package difference

import (
	"fmt"
	"strings"
)

func printMapDiff(diff *Map, indentationLevel int, sign *rune) string {
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
			sign,
			propertyKey,
			propertyValue,
			indentationLevel,
		)

		if hasPreviousPrint && isLastPrint {
			currentSign := rune(currentPrint[0])
			previousSign := rune(previousPrint[0])

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

			if printLine[0] == removed {
				printBuilder.WriteString(printRed(printLine))
			} else if printLine[0] == added {
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
	parentSign *rune,
	propertyKey string,
	propertyValue any,
	indentationLevel int,
) string {
	var propertySign rune
	var propertyPrint strings.Builder

	if parentSign != nil && *parentSign != common {
		propertySign = *parentSign
	} else {
		propertySign = rune(propertyKey[0])
		propertyKey = propertyKey[1:]
	}

	propertyNames := strings.Split(propertyKey, ".")
	propertyName := propertyNames[len(propertyNames)-1]

	prefix, indentation := getPrefixSpacers(propertySign, indentationLevel)

	propertyPrint.WriteString(prefix)
	propertyPrint.WriteString(indentation)
	propertyPrint.WriteString(propertyName)
	propertyPrint.WriteString(": ")

	switch value := (propertyValue).(type) {
	case Map:
		propertyPrint.WriteString("{")
		propertyPrint.WriteString("\n")
	case []Slice:
		propertyPrint.WriteString("[")
		propertyPrint.WriteString("\n")
	default:
		propertyPrint.WriteString(formatValue(prefix+indentation, value))
	}

	switch value := (propertyValue).(type) {
	case Map:
		propertyPrint.WriteString(printMapDiff(&value, indentationLevel+1, &propertySign))
	case []Slice:
		propertyPrint.WriteString(printSliceDiff(&value, indentationLevel+1))
	}

	switch (propertyValue).(type) {
	case Map:
		propertyPrint.WriteString("\n")
		propertyPrint.WriteString(prefix)
		propertyPrint.WriteString(indentation)
		propertyPrint.WriteString("}")
	case []Slice:
		propertyPrint.WriteString("\n")
		propertyPrint.WriteString(prefix)
		propertyPrint.WriteString(indentation)
		propertyPrint.WriteString("]")
	}

	return propertyPrint.String()
}

func printSliceDiff(slices *[]Slice, indentationLevel int) string {
	var printBuilder strings.Builder

	totalCount := len(*slices)

	for index, pair := range *slices {
		if len(pair) != 2 {
			panic("malformed slice pair for slice diff")
		}

		currentSign := pair[0].(rune)
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
			upcomingSign := (*slices)[index+1][0].(rune)
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
