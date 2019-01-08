package response

import (
	"bytes"
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// Model ...
type Model interface {
	Table() string
	Schema() []string
}

// Data ...
type Data struct {
	container *gabs.Container
}

// Container ...
type Container struct {
	v *gabs.Container
}

// To ...
func (c *Container) To(t interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(c.v.EncodeJSON())).Decode(&t)
}

// Path ...
func (d *Data) Path(p string) *Container {
	return &Container{v: d.container.Path(p)}
}

// Container ...
func (d *Data) Container() *gabs.Container {
	return d.container
}

// UnmarshalJSON ...
func (d *Data) UnmarshalJSON(db []byte) error {
	c, err := gabs.ParseJSON(db)
	if err != nil {
		return err
	}
	d.container = c
	return nil
}

// Response ...
type Response struct {
	Data         *Data         `json:"data,omitempty"`
	Errors       []interface{} `json:"errors,omitempty"`
	Subscription chan interface{}
}
