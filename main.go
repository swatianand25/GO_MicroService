package main

import (
	"assign/service"
	"fmt"
	"net/http"
)

func main() {
	//fmt.Println("Hello world")

	//r := service.APIOuter("Vision")

	//fmt.Println(r)

	http.HandleFunc("/", service.APIOuter)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

	fmt.Println("After Handler")
}
