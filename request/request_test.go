package request_test

import (
	"fmt"
	"hasura/request"
	"testing"
)

func TestHasuraRequest(t *testing.T) {
	r := request.HasuraRequest(request.Query("hello"), nil)
	fmt.Println(r.Type())
}
