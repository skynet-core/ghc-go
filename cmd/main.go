package main

import (
	"log"
	"net/http"
	"sync"

	ghc "github.com/skynet-ltd/ghc-go"
)

func main() {

	hc, err := ghc.New("http://0.0.0.0:8080/v1alpha1/graphql", &ghc.Options{
		Header: http.Header{
			"X-Hasura-Access-Key": []string{"SuperPassword"},
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

	var n int
	ids := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		res, err := hc.Execute(ghc.HasuraRequest(ghc.NewMutation(`mutation{
			insert_job(objects:[{
			  op_id:624,
			  job_fields:{},
			  job_result:{},
			  job_type:"sales"
			}]){
			  returning{
				job_id
			  }
			}
		  }`), nil))
		if err != nil {
			log.Fatalln(err)
		}

		if err = res.Data.Path("insert_job.returning.job_id").To(&[]*int{&n}); err != nil {
			log.Fatalln(err)
		}
		ids[i] = n
	}
	wg := sync.WaitGroup{}
	controller := make(chan struct{}, 100)
	for _, id := range ids {
		id := id
		wg.Add(1)
		controller <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-controller
			}()

			if _, err := hc.Execute(ghc.HasuraRequest(ghc.NewMutation(`mutation{
				update_job(where:{job_id:{_eq:%d}},_set:{job_result:{status_code:200}}){
					returning {
					  job_result
					}
				  }
			}`, id), nil)); err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()
	log.Println("DONE")
}
