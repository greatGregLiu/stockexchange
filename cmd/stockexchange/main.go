package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/svett/stockexchange"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", fmt.Sprintf(":%d", GetenvInt("PORT", 9292)), "Listen and serve HTTP on host:port")
}

func main() {
	flag.Parse()

	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("The provided %s addr is not correct format (ex. host:port)", addr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", stockexchange.Index)

	server := negroni.Classic()
	server.UseHandler(router)

	log.Printf("StackExchange started. HTTP listen and serve on %s", addr)
	http.ListenAndServe(addr, server)
}

func GetenvInt(name string, defaultValue int) int {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			return value
		}
	}
	return defaultValue
}