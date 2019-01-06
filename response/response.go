package response

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type data struct {
}

// Response ...
type Response struct {
	Data         map[string][]map[string]interface{} `json:"data,omitempty"`
	Errors       []interface{}                       `json:"errors,omitempty"`
	Subscription chan interface{}
}

// ConsvertTo ...
func (res *Response) ConsvertTo(sv interface{}) (err error) {
	defer func() {
		if ev := recover(); ev != nil {
			err = fmt.Errorf("%v", ev)
		}
	}()

	t := reflect.TypeOf(sv).Elem().Elem()
	result := reflect.New(reflect.SliceOf(t))

	for _, sv := range res.Data {
		for _, v := range sv {
			nt := reflect.New(t)
			if err = mapstructure.Decode(v, nt.Interface()); err != nil {
				return err
			}
			result.Elem().Set(reflect.Append(result.Elem(), nt.Elem()))
		}
	}

	reflect.ValueOf(sv).Elem().Set(result.Elem())
	return nil
}
