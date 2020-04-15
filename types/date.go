package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

//Date é a data sem o tempo. Converte em yyyy-mm-dd
type Date time.Time

//UnmarshalJSON implementação
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

//MarshalJSON implementação
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

//Value implementação
func (d Date) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

//NewDate cria um novo Date a partir de uma string yyyy-mm-dd
func NewDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date(t), nil
}
