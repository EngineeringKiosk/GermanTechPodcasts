package io

import (
	"encoding/json"
	"io/ioutil"
)

// WriteJSONFile will write v into absFilePath.
// v needs to be JSON marshalable.
func WriteJSONFile(absFilePath string, v any) error {
	content, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(absFilePath, content, 0644)
	return err
}
