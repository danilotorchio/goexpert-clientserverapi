package main

import (
	"log"
	"net/http"

	"github.com/danilotorchio/goexpert-clientserverapi/database"
	"github.com/danilotorchio/goexpert-clientserverapi/handler"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	http.HandleFunc("/cotacao", handler.ExchangeHandler)
	http.HandleFunc("/historico", handler.HistoryHandler) // Bonus

	log.Default().Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
