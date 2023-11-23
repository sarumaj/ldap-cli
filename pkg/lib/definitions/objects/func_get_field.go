package objects

import "reflect"

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
