package utils

import (
	"html"
	"reflect"
)

func ScapeHTML(value reflect.Value){

	// loop over the struct
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		// check if the field is a string
		if field.Type() != reflect.TypeOf("") {
			continue
		}

		str := field.Interface().(string)
		// set field to escaped version of the string
		field.SetString(html.EscapeString(str))
	}
}
