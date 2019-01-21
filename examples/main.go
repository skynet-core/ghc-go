package main

import (
	"fmt"
	"log"
	"strings"

	ghc "github.com/skynet-ltd/ghc-go"
)

// Sales ...
type Sales struct {
	Geo     string
	Keyword string
}

// Schema ...
func (s *Sales) Schema() []string {
	return []string{"geo", "keyword"}
}

// Table ...
func (s *Sales) Table() string {
	return "sales"
}

func main() {

	c, err := ghc.New("http://localhost:8080/v1alpha1/graphql", nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := c.Execute(ghc.HasuraRequest(
		ghc.NewQuery(`query{
					sales(where:{id:{_lte:%d}}){ %s }
				}`, 3, strings.Join([]string{"geo", "keyword"}, ",")), nil,
	))
	if err != nil {
		log.Fatalln(err)
	}
	geos := []string{}
	if err = resp.Data.Path("sales.keyword").To(&geos); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(geos[0])

}
