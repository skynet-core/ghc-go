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

// Schema ...
func (s *Sales) Schema() []string {
	return []string{"geo", "keyword"}
}

// Table ...
func (s *Sales) Table() string {
	return "sales"
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
	if err = resp.MapResult(&sales); err != nil {
		log.Fatalln(err)
	}

	for _, s := range sales {
		fmt.Println(s)
	}

	resp, err = c.Execute(request.HasuraRequest(
		request.NewMutation(`mutation ($jobs:[job_insert_input!]!){
					insert_job(objects:$jobs){
						affected_rows
						returning {
							job_id
						}
					}
				}`), request.Variables{
			"jobs": []map[string]interface{}{
				map[string]interface{}{
					"op_id":      244,
					"job_type":   "sales",
					"job_fields": map[string]interface{}{},
				},
			},
		},
	))
	if err != nil {
		log.Fatalln(err)
	}

	dst := []map[string]interface{}{}

	if err = resp.Returning(&dst); err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.AffectedRows())

	resp, err = c.Execute(request.HasuraRequest(
		request.Query(`query{
					job_aggregate(where: {job_id:{_gte:0}}){
						aggregate {
							count
						}
					}
				}`), nil,
	))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.Aggregate())

}
