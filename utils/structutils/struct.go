package structutils

import (
	"reflect"
	"strings"
)

func GetFieldTagValue(structPointer interface{}, fieldPointer interface{}) string {

	var tagValue string

	structReflect := reflect.ValueOf(structPointer).Elem()
	fieldReflect := reflect.ValueOf(fieldPointer).Elem()

	for i := 0; i < structReflect.NumField(); i++ {
		fieldValue := structReflect.Field(i)
		if fieldValue.Addr().Interface() == fieldReflect.Addr().Interface() {
			tagValue = structReflect.Type().Field(i).Tag.Get("mapstructure")
			if tagValue == "" {
				tagValue = structReflect.Type().Field(i).Tag.Get("yaml")
			}
			if tagValue == "" {
				tagValue = structReflect.Type().Field(i).Tag.Get("json")
			}
		}
	}

	return strings.Split(tagValue, ",")[0]
}
