package testdata

import (
	"database/sql"
	"sync"
)

type MyStruct struct {
	B int
	A int `json:"a"`
	C int
}

type MyStruct2 struct {
	X int
	Y int
	Z int

	V string
	C int
}

type Address struct {
	Address2 sql.NullString `json:"address_2"`
	Address1 sql.NullString
	UserID   sql.NullInt64 `json:"user_id"`

	OK *bool
}

type EmbedStruct struct {
	B struct {
		B struct{}
		A struct{}
	}
	// A
	A struct {
		B struct{}
		A struct{}
	}
}

type StructWithComment struct {
	// C
	C string

	// G
	G string

	// A
	B string
	A string
}

type StructWithMutex struct {
	B string

	mu sync.Mutex
	C  string
}
