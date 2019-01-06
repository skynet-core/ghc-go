package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/skynet-ltd/ghc-go/client"
	"github.com/skynet-ltd/ghc-go/request"
)

// Sales ...
type Sales struct {
	Geo     string `mapstructure:"geo"`
	Keyword string `mapstructure:"keyword"`
}

func main() {
	c, err := client.New("http://localhost:8080/v1alpha1/graphql", nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := c.Execute(request.HasuraRequest(
		request.NewQuery(`query{
					sales(where:{id:{_lte:%d}}){ %s }
				}`, 3, strings.Join([]string{"geo", "keyword"}, ",")), nil,
	))
	if err != nil {
		log.Fatalln(err)
	}

	sales := make([]*Sales, 0)
	if err = resp.ConsvertTo(&sales); err != nil {
		log.Fatalln(err)
	}

	for _, s := range sales {
		fmt.Println(*s)
	}

}
