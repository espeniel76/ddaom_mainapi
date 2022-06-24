package main

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/memdb"
	"ddaom/mlogdb"
	"ddaom/rest"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("ddaom API v0.3")
	setInitialize()
}

func setInitialize() {
	define.SetDefineApiParse()
	db.RunMySql()
	memdb.RunRedis()
	mlogdb.RunMongodb()

	mux := mux.NewRouter()
	rest.Handlers(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	// for REST HTTP
	err := http.ListenAndServe(":"+define.HTTP_PORT, handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	// for REST SSL
	// go func() {
	// err := http.ListenAndServeTLS(":"+define.HTTP_PORT_SSL, "/usr/local/ssl/cert.pem", "/usr/local/ssl/key.pem", handler)
	// if err != nil {
	// 	log.Fatal("ListenAndServeTLS:", err)
	// }
	// }()

	fmt.Println("Test Edit")
}
