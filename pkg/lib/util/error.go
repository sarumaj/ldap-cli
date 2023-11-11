package util

func PanicIfError[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
