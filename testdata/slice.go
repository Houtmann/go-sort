package testdata

var mymap = map[string]string{
	"key":  "value",
	"data": "data",
}

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

const (
	Waiting int64 = iota
	Succeeded
	Failed
)

const (
	WaitingStr   = "waiting"
	SucceededStr = "succeeded"
	FailedStr    = "failed"
)

var Status = map[int64]string{
	Waiting:   WaitingStr,
	Succeeded: SucceededStr,
	Failed:    FailedStr,
}
