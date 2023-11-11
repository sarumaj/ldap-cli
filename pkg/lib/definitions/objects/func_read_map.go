package objects

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func lookupMapKey(key string, raw map[string]any) (value any, ok bool) {
	value, ok = raw[key]
	if ok {
		return
	}

	for k, v := range raw {
		if strings.EqualFold(key, k) {
			return v, true
		}
	}

	return
}

func readMap[T any](o *T, raw map[string]any) error {
	v := reflect.ValueOf(o).Elem()

	var errs []error
	for j := 0; j < v.NumField(); j++ {
		field := v.Field(j)

		// lookup attribute name by tag
		var map_key string
		if tag := v.Type().Field(j).Tag.Get("ldap_attr"); tag != "" {
			map_key = tag
		} else {
			map_key = v.Type().Field(j).Name
		}

		if map_key == "-" {
			continue
		}

		if value, ok := lookupMapKey(map_key, raw); ok && field.IsValid() {
			func() {
				// handle panic
				defer func() {
					if recovered := recover(); recovered != nil {
						errs = append(errs, fmt.Errorf("%v", recovered))
					}
				}()

				// type handling
				// new types can be added here
				switch field.Kind() {

				case reflect.Array, reflect.Slice:
					switch field.Type().Elem().Kind() {

					case reflect.String:
						if _v, ok := value.(string); ok {
							value = []string{_v}
						}

					}

				case reflect.Bool:
					value = util.PanicIfError(strconv.ParseBool(value.(string)))

				case reflect.Int:
					value = util.PanicIfError(strconv.Atoi(value.(string)))

				case reflect.Int64:
					value = util.PanicIfError(strconv.ParseInt(value.(string), 10, 64))

				case reflect.String:
					if _v, ok := value.([]string); ok {
						value = strings.Join(_v, ";")
					}

				}

				field.Set(reflect.ValueOf(value))
			}()
		}
	}

	return errors.Join(errs...)
}
