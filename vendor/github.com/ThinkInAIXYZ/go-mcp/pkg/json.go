package pkg

import (
	"encoding/json"
	"fmt"
)

// var sonicAPI = sonic.Config{UseInt64: true}.Froze() // Effectively prevents integer overflow

func JSONUnmarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%w: data=%s, error: %+v", ErrJSONUnmarshal, data, err)
	}
	return nil
}
