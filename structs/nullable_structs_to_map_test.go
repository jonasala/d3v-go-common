package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
)

type testStruct struct {
	Field1 string      `db:"-"`
	Field2 null.String `db:"field2"`
	Field3 null.Int    `db:"field3"`
}

func TestNullableStructToMap(t *testing.T) {
	s := testStruct{
		Field1: "Test",
		Field2: null.StringFrom("Valid String"),
		Field3: null.Int{Int: 0, Valid: false},
	}

	m, err := NullableStructToMap(s)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Valid String", m["field2"])
}
