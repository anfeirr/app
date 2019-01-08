package core

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// JsFormat formats a javascript expression according to a format specifier.
func JsFormat(format string, a ...interface{}) (string, error) {
	for i, v := range a {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array,
			reflect.String:
			b, err := json.Marshal(v)
			if err != nil {
				return "", errors.Wrapf(err, "converting %T to json failed", v)
			}

			a[i] = string(b)

		case reflect.Func:
			return "", errors.Errorf("formatting funcs is not supported: %T", v)
		}
	}

	return fmt.Sprintf(format, a...), nil
}
