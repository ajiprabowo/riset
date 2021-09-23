package main

import (
	"embed"

	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/ajiprabowo/riset/statik"
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

// ref. https://golang.hotexamples.com/examples/http/-/ServeFile/golang-servefile-function-examples.html
func readme(w http.ResponseWriter, r *http.Request) {

	//Using helper ServeFile
	//http.ServeFile(w, r, "readme.txt")

	//manual
	//body, err := ioutil.ReadFile("readme.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Fprintf(w, "%s", body)

	//statik embed
	//statikFS, err := fs.New()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Access individual files by their paths.
	//f, err := statikFS.Open("/readme.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()

	//contents, err := ioutil.ReadAll(f)

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
