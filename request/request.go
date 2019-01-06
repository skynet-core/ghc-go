package request

import (
	"fmt"
	"reflect"
)

// HasuraQuery ...
type HasuraQuery interface {
	String() string
}

// Query ...
type Query string

func (q Query) String() string {
	return string(q)
}

// Mutation  ....
type Mutation string

func (q Mutation) String() string {
	return string(q)
}

// Subscription ...
type Subscription string

func (q Subscription) String() string {
	return string(q)
}

// NewQuery ...
func NewQuery(str string, args ...interface{}) HasuraQuery {
	return Query(fmt.Sprintf(str, args...))
}

// NewMutation ...
func NewMutation(str string, args ...interface{}) HasuraQuery {
	return Mutation(fmt.Sprintf(str, args...))
}

// NewSubscription ...
func NewSubscription(str string, args ...interface{}) HasuraQuery {
	return Subscription(fmt.Sprintf(str, args...))
}

// Variables ...
type Variables map[string]interface{}

// Request ...
type Request struct {
	Query        `json:"query,omitempty"`
	Mutation     `json:"mutation,omitempty"`
	Subscription `json:"subscription,omitempty"`
	Variables    `json:"variables,omitempty"`

	t reflect.Type
}

// Type ...
func (r *Request) Type() reflect.Type {
	return r.t
}

// HasuraRequest ...
func HasuraRequest(query HasuraQuery, vars Variables) *Request {
	rq := &Request{
		t: reflect.TypeOf(query),
	}
	switch v := query.(type) {
	case Query:
		rq.Query = Query(query.String())
	case Mutation:
		rq.Mutation = Mutation(query.String())
	case Subscription:
		rq.Subscription = Subscription(query.String())
	default:
		panic("there are no registered cases for " + v.String())
	}

	rq.Variables = vars
	return rq
}
