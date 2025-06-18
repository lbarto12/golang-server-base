package main

import (
	"fmt"
	"golang-server-base/api"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()

	server := api.Server{
		Host: "localhost",
		Port: "8080",
		Cors: &cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost},
			AllowCredentials: true,
		},
	}

	server.AddHandlers(map[string]http.Handler{
		"/ping": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}),
		"/yurr": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("skrt"))
		}),
	})

	go Test()

	log.Fatal(server.Launch())
}

func Test() {
	// time.Sleep(1 * time.Second)
	fmt.Println("Tests Starting...")
	// time.Sleep(1 * time.Second)

	// err := minioapi.Init(minioapi.EnvGetOptions())
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println("Connected to MINIO")

	// client := minioapi.Client()

	// err = client.MakeBucket(context.TODO(), "test-bucket", minio.MakeBucketOptions{
	// 	Region: "us-east-1",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.RemoveBucket(context.TODO(), "test-bucket")
	// if err != nil {
	// 	panic(err)
	// }

}
