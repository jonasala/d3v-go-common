package squirrelhelpers

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jonasala/d3v-go-common/structs"
)

//UpdateFields configures a squirrel.UpdateBuilder appending update fields
func UpdateFields(builder squirrel.UpdateBuilder, data interface{}, timestamps ...string) (squirrel.UpdateBuilder, error) {
	fields, err := structs.NullableStructToMap(data)
	if err != nil {
		return builder, err
	}

	now := time.Now()
	for _, field := range timestamps {
		fields[field] = now
	}

	for field, value := range fields {
		builder = builder.Set(field, value)
	}
	return builder, nil
}
