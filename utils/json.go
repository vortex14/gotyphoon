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
