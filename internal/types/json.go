package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON represents a JSON field type for SQLBoiler
type JSON json.RawMessage

// Value implements the driver.Valuer interface for JSON
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

// Scan implements the sql.Scanner interface for JSON
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
	case string:
		*j = JSON(v)
	default:
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for JSON
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for JSON
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("JSON: UnmarshalJSON on nil pointer")
	}
	*j = JSON(data)
	return nil
}

// String returns the string representation of the JSON field
func (j JSON) String() string {
	return string(j)
}
