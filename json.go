// Copyright (c) 2024 Daniel van Dijk (https://daniel.vandijk.sh)
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package difference

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
