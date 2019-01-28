package main

import (
	"context"
	"log"
	"net/http"
	"time"

	ghc "github.com/skynet-ltd/ghc-go"
)

func main() {
	result := make(chan time.Time, 0)
	mux := http.NewServeMux()

	mux.HandleFunc("/results", func(w http.ResponseWriter, req *http.Request) {
		result <- time.Now()
		w.WriteHeader(200)
	})

	srv := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	defer func() {
		srv.Shutdown(context.Background())
	}()

	go func() {
		srv.ListenAndServe()
	}()

	hc, err := ghc.New("http://0.0.0.0:8080/v1alpha1/graphql", &ghc.Options{
		Header: http.Header{
			"X-Hasura-Access-Key": []string{"SuperPassword"},
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

	start := time.Now()

	_, err = hc.Execute(ghc.HasuraRequest(
		ghc.NewMutation(`mutation{
			update_job(where:{job_id:{_eq:11280}},_inc:{tries:1},_set:{job_result:"{\"test\":30}"}){
				affected_rows
			}
		}`), nil,
	))

	if err != nil {
		log.Println(err)
		return
	}

	end := <-result
	log.Println("Test finished: ", end.UnixNano()-start.UnixNano())
}
