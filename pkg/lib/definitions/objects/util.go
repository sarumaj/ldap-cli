package objects

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sarumaj/ldap-cli/pkg/lib/util"
)

func GetDistinguishedNames[T any](objects []T) []string {
	return GetProperties[string](objects, "DistinguishedName")
}

func GetField[S, T any](o *T, property string) S {
	to := reflect.New(reflect.TypeOf((*S)(nil)).Elem())
	toValue := reflect.Indirect(to)
	field := reflect.Indirect(reflect.ValueOf(o)).FieldByName(property)

	if toValue.CanConvert(field.Type()) {
		toValue.Set(toValue.Convert(field.Type()))
	}

	if toValue.CanSet() && field.Type().AssignableTo(toValue.Type()) {
		toValue.Set(field)
	}

	return toValue.Interface().(S)
}

func GetProperties[S, T any](objects []T, property string) (properties []S) {
	for _, o := range objects {
		properties = append(properties, GetField[S](&o, property))
	}

	return
}

func hexify(s string) string {
	var chars []string
	for _, c := range []byte(s) {
		chars = append(chars, fmt.Sprintf("\\x%02X", c))
	}

	return strings.Join(chars, "")
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

		if value, ok := raw[strings.ToLower(map_key)]; ok && field.IsValid() {

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
		}
	}

	return errors.Join(errs...)
}

func Unhexify(s string) string {
	chars := []byte{}
	for _, c := range strings.Split(strings.TrimLeft(s, "\\x"), "\\x") {
		i, _ := strconv.ParseInt(c, 16, 8)
		chars = append(chars, byte(i))
	}
	return string(chars)
}
