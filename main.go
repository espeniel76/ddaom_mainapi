package main

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/memdb"
	"ddaom/rest"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("ddaom API v0.4")
	setInitialize()
}

func setInitialize() {
	define.SetDefineApiParse()
	define.SetConnectionInfosParse()
	db.RunMySql()
	memdb.RunRedis()
	// mlogdb.RunMongodb()

	mux := mux.NewRouter()
	rest.Handlers(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	// for REST HTTP
	err := http.ListenAndServe(":"+define.Mconn.HTTPPort, handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
