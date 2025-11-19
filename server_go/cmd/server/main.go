package main

import (
	"fmt"
	"net/http"
	"server/internal/account"
	"server/internal/handler"
)

func main(){
	// init accounts
	account.NewAccount("1", "Max", 1000.0)
	account.NewAccount("2", "John", 2000.0)
	
	// setup routes
	http.HandleFunc("/balance", handler.Balance)
	
	// start server
	port := ":8080"
	fmt.Println("Server running on", port)
	http.ListenAndServe(port, nil)
}