package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

//Date é a data sem o tempo. Converte em yyyy-mm-dd
type Date time.Time

//UnmarshalJSON implementação
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		*d = Date(t)
	}
	return nil
}

//MarshalJSON implementação
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(t.Format("2006-01-02"))
}

//Value implementação
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format("2006-01-02"), nil
}

//Scan implementação
func (d *Date) Scan(value interface{}) error {
	if value != nil {
		switch value.(type) {
		case time.Time:
			*d = Date(value.(time.Time))
		case string:
			t, err := time.Parse("2006-01-02", value.(string))
			if err != nil {
				return err
			}
			*d = Date(t)
		default:
			return errors.New("invalid date field")
		}
	}
	return nil
}

//NewDate cria um novo Date a partir de uma string yyyy-mm-dd
func NewDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date(t), nil
}
