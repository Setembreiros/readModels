package common_data

import "encoding/json"

func DeserializeData(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}
