package response

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var returnedOps = map[string]bool{
	"insert": true,
	"update": true,
	"delete": true,
}

// Model ...
type Model interface {
	Table() string
	Schema() []string
}

// Response ...
type Response struct {
	Data         `json:"data,omitempty"`
	Errors       []interface{} `json:"errors,omitempty"`
	Subscription chan interface{}
}

// Data ...
type Data map[string]json.RawMessage

// Returning ...
func (d Data) Returning(dst interface{}) error {
	for k, v := range d {
		if ss := strings.Split(k, "_"); len(ss) > 1 && returnedOps[ss[0]] {
			m := make(map[string]json.RawMessage, 0)
			if err := json.Unmarshal(v, &m); err != nil {
				return err
			}
			return json.Unmarshal(m["returning"], dst)
		}
	}
	return errors.New("no returning objects")
}

// MapResult ...
func (d Data) MapResult(m interface{}) error {

	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Slice {
		return errors.New("argument must have *[]response.Model type")
	}

	inst := reflect.New(t.Elem().Elem()).Elem().Interface()
	ni, ok := inst.(Model)
	if !ok {
		return errors.New("argument must have *[]response.Model type")
	}
	if r, ok := d[ni.Table()]; ok {
		return json.Unmarshal(r, m)
	}
	return nil
}

// AffectedRows ...
func (d Data) AffectedRows() (int, error) {
	for k, v := range d {
		if ss := strings.Split(k, "_"); len(ss) > 1 && returnedOps[ss[0]] {
			m := make(map[string]json.RawMessage, 0)
			if err := json.Unmarshal(v, &m); err != nil {
				return 0, err
			}
			var rows int
			return rows, json.Unmarshal(m["affected_rows"], &rows)
		}
	}
	return 0, errors.New("no returning objects")
}

// Aggregate ...
func (d Data) Aggregate() map[string]interface{} {
	for k, v := range d {
		if ss := strings.Split(k, "_"); len(ss) > 1 && ss[1] == "aggregate" {
			m := make(map[string]interface{}, 0)
			if err := json.Unmarshal(v, &m); err != nil {
				return nil
			}
			return m
		}
	}
	return nil
}
