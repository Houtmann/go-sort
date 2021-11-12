package testdata

type MyStruct struct {
	A int
	B int
	C int
}

type MyStruct2 struct {
	D int
	E string
	X int `json:"x"`
	Y int
	Z int `json:"z"`
}

type Address struct {
	Address1 string
	Address2 string
	City     string
	Phone    string
	UserID   int64
}

//
// func MyTestFunc() {
// 	address := Address{
// 		Address1: "address1",
// 		Address2: "address2",
//
// 		Phone: "+33XXXXXXXX",
// 		City:  "NY",
// 	}
// 	fmt.Println(address)
// 	_ = struct {
// 		BA string
// 		BC string
// 	}{
// 		BC: "",
// 		BA: "",
// 	}
// }
