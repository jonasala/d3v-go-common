package squirrelhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Masterminds/squirrel"
	"github.com/volatiletech/null"
)

type testStruct struct {
	Field1 string      `db:"-"`
	Field2 null.String `db:"field2"`
	Field3 null.Int    `db:"field3"`
}

func TestInsertFields(t *testing.T) {
	builder := squirrel.Insert("test_table")
	s := testStruct{
		Field1: "test1",
		Field2: null.StringFrom("test2"),
		Field3: null.Int{Int: 0, Valid: false},
	}
	builder, err := InsertFields(builder, s, "created_at")
	assert.NoError(t, err)

	query, _, err := builder.ToSql()

	expectedSQL :=
		"INSERT INTO test_table (field2,created_at) VALUES (?,?)"

	assert.NoError(t, err)
	assert.Equal(t, expectedSQL, query)
}
