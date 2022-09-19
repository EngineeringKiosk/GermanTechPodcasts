package io

import (
	"encoding/json"
	"os"
)

// WriteJSONFile will write v into absFilePath.
// v needs to be JSON marshalable.
func WriteJSONFile(absFilePath string, v any) error {
	content, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(absFilePath, content, 0644)
	return err
}
