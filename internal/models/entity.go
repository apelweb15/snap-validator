package models

import "reflect"

type CustomValidator struct {
	FieldType  reflect.StructField
	FieldValue reflect.Value
}

type ValidatorProperty struct {
	FieldName  string
	JsonName   string
	Array      bool
	FieldType  reflect.StructField
	FieldValue reflect.Value
	IdxArray   int
}
