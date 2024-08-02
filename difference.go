package difference

import (
	"io"
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
