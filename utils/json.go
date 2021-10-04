package utils

import (
	"encoding/json"
)

func JsonLoad(model interface{}, data string) error  {
	if err := json.Unmarshal([]byte(data), model); err != nil {
		return err
	}
	return nil
}

func JsonDump(model interface{}) (error, []byte) {
	data, err := json.Marshal(model); if err != nil { return  err, nil }
	return nil, data
}

func PrintPrettyJson(f interface{}) (error, string) {
	dump, err := json.MarshalIndent(f, "  ", "  ")
	return err, string(dump)
}