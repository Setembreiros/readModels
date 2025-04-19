package test_common

import "encoding/json"

func SerializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}
