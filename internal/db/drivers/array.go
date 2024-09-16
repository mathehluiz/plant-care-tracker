package drivers

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(a, ",") + "}", nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		fv := string(v)[1:]
		fv = fv[:len(fv)-1]
		if fv == "" {
			*a = make(StringArray, 0)
			return nil
		}

		*a = strings.Split(fv, ",")
		return nil
	default:
		return errors.New("unsupported Scan, storing driver.Value type " + v.(string) + " into type *[]string")
	}
}
