package testdata

import (
	"database/sql"
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

	HasShippings *bool
}
