package main

import (
	"fmt"
	"job-system/controller/handler"
	"job-system/repository"
	"job-system/service"
	"log"
	"net/http"
)

func main() {
	// application logic here
	/*Dependency Injection
	main()
	 │
	 ├─ creates Repository
	 │
	 ├─ gives Repository to Service
	 │
	 └─ gives Service to Handler/
	 repo := NewJobRepository()
	 servie := NewJobService(repo)
	 handler := NewJobHandler(service)
	*/
	repo := repository.NewJobRepository()
	srv := service.NewJobService(repo)
	jobHandler := handler.NewJobHandler(srv)

	http.DefaultServeMux.HandleFunc("POST /jobs", jobHandler.CreateJob)
	http.DefaultServeMux.HandleFunc("GET /jobs/{id}", jobHandler.GetJob)
	http.DefaultServeMux.HandleFunc("DELETE /jobs/{id}", jobHandler.Delete)
	fmt.Printf("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
