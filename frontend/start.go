package main

import (
	"go-clawer/frontend/controller"
	"net/http"
)

func main() {
	http.Handle("/search", controller.CreateSearchResultHandler("frontend/view/template.html"))
	http.Handle("/", http.FileServer(http.Dir("frontend/view")))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
