package main

import (
	"embed"

	"fmt"
	"log"
	"net/http"
	"time"
)

//go:embed readme.txt
var f embed.FS

func headers(w http.ResponseWriter, r *http.Request) {

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func context(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	fmt.Println("server: context handler started")
	defer fmt.Println("server: context handler ended")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "hello context\n")
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func readme(w http.ResponseWriter, r *http.Request) {

	contents, err := f.ReadFile("readme.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", contents)
}

func main() {

	http.HandleFunc("/headers", headers)
	http.HandleFunc("/context", context)
	http.HandleFunc("/readme", readme)
	http.ListenAndServe(":8090", nil)
}
