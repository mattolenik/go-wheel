package wheel

import (
	"encoding/json"
	"fmt"
)

// JSON is a type alias of map[string]any for use with JSON arguments at the command line.
type JSON map[string]any

func (j JSON) String() string {
	return j.ToString()
}

func (j JSON) ToString() string {
	data, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("cannot display value due to error: %w", err).Error()
	}
	return string(data)
}

func (j JSON) FromString(s string) error {
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return fmt.Errorf("cannot parse value due to error: %w", err)
	}
	return nil
}
