package main

import (
	awss3 "microservices-aws/pkg/aws-s3"
)

func main() {
	awss3.Example()
	// router := mux.NewRouter()
	// routes.SetRoutes(router)
	// srv := &http.Server{
	// 	Handler:      router,
	// 	Addr:         "0.0.0.0:8001",
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }
	// fmt.Println("Server 8001")
	//log.Fatal(srv.ListenAndServe())
}
