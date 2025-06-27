package main

import (
	"fmt"
	"net/http"
)

func main(){
	fmt.Println("GoTasker Pro API starting...")

	http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		fmt.Fprintln(w, "Hello from GoTasker Pro!")
	})

	http.ListenAndServe(":8080",nil)
}