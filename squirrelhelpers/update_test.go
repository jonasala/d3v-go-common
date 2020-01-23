package squirrelhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Masterminds/squirrel"
	"github.com/volatiletech/null"
)

func TestUpdateFields(t *testing.T) {
	builder := squirrel.Update("test_table")
	s := testStruct{
		Field1: "test1",
		Field2: null.StringFrom("test2"),
		Field3: null.Int{Int: 0, Valid: false},
	}
	builder, err := UpdateFields(builder, s, "updated_at")
	assert.NoError(t, err)

	query, _, err := builder.ToSql()

	expectedSQL :=
		"UPDATE test_table SET field2 = ?, updated_at = ?"

	assert.NoError(t, err)
	assert.Equal(t, expectedSQL, query)
}
