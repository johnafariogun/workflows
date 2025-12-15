package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{
			"Laptop", "Keyboard", "Phone",
		})
	})

	log.Println("Inventory service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
