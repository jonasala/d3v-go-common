package structs

import (
	"errors"
	"reflect"
	"strings"
)

//NullableStructToMap converts a struct containing nullable fields (https://github.com/volatiletech/null) tagged with `db:"somefield"` in a map of valid fields
func NullableStructToMap(s interface{}) (map[string]interface{}, error) {

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("s is not a struct")
	}

	validFields := map[string]interface{}{}
	loadFields(v, &validFields)

	return validFields, nil
}

func loadFields(v reflect.Value, fields *map[string]interface{}) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)

		if tag, _ := f.Tag.Lookup("db"); tag != "" && tag != "-" {
			if strings.HasPrefix(f.Type.String(), "null.") {
				if valid := v.Field(i).FieldByName("Valid").Bool(); valid {
					(*fields)[tag] = v.Field(i).Field(0).Interface()
				}
			} else {
				(*fields)[tag] = v.Field(i).Interface()
			}
		} else if tag != "-" {
			nf := v.Field(i)
			if nf.Kind() == reflect.Ptr {
				nf = nf.Elem()
			}

			if nf.Kind() == reflect.Struct {
				loadFields(nf, fields)
			}
		}
	}
}
