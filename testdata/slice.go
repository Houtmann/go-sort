package testdata

func MyFunc() ([]string, []string) {
	newslice := []string{"boat", "plane", "bike"}
	newslice2 := []string{"boat",
		"plane",
		"bike",
	}
	_ = map[string]string{
		"b": "value",
		"c": "data",
	}

	return newslice, newslice2
}
