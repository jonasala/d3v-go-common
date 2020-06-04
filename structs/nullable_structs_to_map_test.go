package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
)

type testStruct struct {
	Field1 *string     `db:"-"`
	Field2 null.String `db:"field2"`
	Field3 null.Int    `db:"field3"`
	*testNestedStruct
}

type testNestedStruct struct {
	Field4 int         `db:"field4"`
	Field5 null.String `db:"field5"`
}

func TestNullableStructToMap(t *testing.T) {
	s := testStruct{
		Field1: nil,
		Field2: null.StringFrom("Valid String"),
		Field3: null.Int{Int: 0, Valid: false},
		testNestedStruct: &testNestedStruct{
			Field4: 2,
			Field5: null.StringFrom("Valid nested string"),
		},
	}

	m, err := NullableStructToMap(s)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Valid String", m["field2"])
	assert.Equal(t, 2, m["field4"])
	assert.Equal(t, "Valid nested string", m["field5"])
}
