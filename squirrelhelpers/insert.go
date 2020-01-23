package squirrelhelpers

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jonasala/d3v-go-common/structs"
)

//InsertFields configure a squirrel.InsertBuilder appending insert fields
func InsertFields(builder squirrel.InsertBuilder, data interface{}, timestamps ...string) (squirrel.InsertBuilder, error) {
	fields, err := structs.NullableStructToMap(data)
	if err != nil {
		return builder, err
	}

	now := time.Now()
	for _, field := range timestamps {
		fields[field] = now
	}

	columns := []string{}
	values := []interface{}{}
	for field, value := range fields {
		columns = append(columns, field)
		values = append(values, value)
	}
	return builder.Columns(columns...).Values(values...), nil
}
