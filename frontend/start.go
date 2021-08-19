package main

import (
	"go-clawer/frontend/controller"
	"net/http"
)

func main() {
	http.Handle("/search", controller.CreateSearchResultHandler("frontend/view/template.html"))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
