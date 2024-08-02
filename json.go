package indifferent

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrJsonDecode = errors.New("failed to decode json")

func decodeJsonIntoMap(reader io.Reader, buffer *Map) error {
	if err := json.NewDecoder(reader).Decode(buffer); err != nil {
		return ErrJsonDecode
	} else {
		return nil
	}
}
