package testdata

func VariadicFunc(keys ...string) []string {
	return keys
}

func NoneVariadicFunc(id string, password string) []string {
	return nil
}

func checkVariadic() {
	VariadicFunc(
		"B",
		"C",
		"A")
	NoneVariadicFunc("t", "r")
}
