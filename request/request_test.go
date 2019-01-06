package request_test

import (
	"fmt"
	"testing"

	"github.com/skynet-ltd/ghc-go/request"
)

func TestHasuraRequest(t *testing.T) {
	r := request.HasuraRequest(request.Query("hello"), nil)
	fmt.Println(r.Type())
}
