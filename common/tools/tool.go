package tools

import (
	"bytes"
	"encoding/json"
	"os"
)

func IndentFromBody(data []byte) error {
	var output bytes.Buffer
	err := json.Indent(&output, data, "", "\t")
	if err != nil {
		return err
	}

	output.WriteTo(os.Stdout)
	return nil
}
