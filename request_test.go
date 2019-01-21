package ghc_test

import (
	"fmt"
	"testing"

	ghc "github.com/skynet-ltd/ghc-go"
)

func TestHasuraRequest(t *testing.T) {
	r := ghc.HasuraRequest(ghc.Query("hello"), nil)
	fmt.Println(r.Type())
}
