package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Keritos/tesseract"
)

var (
	ocrpath string
)

func init() {
	flag.StringVar(&ocrpath, "p", "", `path\to\tesseract.exe`)
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func cap2txt(w http.ResponseWriter, r *http.Request) error {
	txt, err := tesseract.ReadText(r.Body)
	if err != nil {
		return err
	}
	fmt.Fprint(w, txt)
	return nil
}

func main() {
	flag.Parse()

	if ocrpath == "" {
		panic("tesseract.exe path not found")
	}

	tesseract.ExecutablePath = ocrpath
	http.Handle("/cap2txt", appHandler(cap2txt))
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { log.Println("pong!!!") })
	http.ListenAndServe(":"+os.Getenv("HTTP_PLATFORM_PORT"), nil)
}
