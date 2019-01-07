package response

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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

// AffectedRows ...
func (d Data) AffectedRows() int {
	n, ok := d["affected_rows"]
	if !ok {
		return 0
	}
	var num int
	if err := json.Unmarshal(n, &num); err != nil {
		return 0
	}
	return num
}

// ReturningRaw ...
func (d Data) ReturningRaw() json.RawMessage {
	if rm, ok := d["returning"]; ok {
		return rm
	}
	return nil
}

// ReturningMap ...
func (d Data) ReturningMap() map[string]interface{} {
	if rm, ok := d["returning"]; ok {
		m := make(map[string]interface{})
		if err := json.Unmarshal(rm, &m); err != nil {
			return nil
		}
		return m
	}
	return nil
}

// MapResult ...
func (d Data) MapResult(m *[]Model) (err error) {
	defer func() {
		if ev := recover(); ev != nil {
			err = fmt.Errorf("%v", ev)
		}
	}()

	inst := reflect.New(reflect.TypeOf(m).Elem().Elem()).Elem().Interface()
	ni := inst.(Model)
	// if !ok {

	// }
	if r, ok := d[ni.Table()]; ok {
		err = json.Unmarshal(r, m)
		return
	}
	return
}
