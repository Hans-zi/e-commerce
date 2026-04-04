package utils

import "encoding/json"

func Copy(dest interface{}, src interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	return nil
}
