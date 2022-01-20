package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Produce: structured list of produce items in market
type Produce struct {
	ProduceCode string `json:"code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"price"`
}

// Food: Array of Produce items
var Food []Produce

// GetAllFoods: Displays all items in produce list
func GetAllFoods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /produce: GetAllFoods")
	json.NewEncoder(w).Encode(Food)
}

func GetFoodByCode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /produce/{code}: GetFoodByCode")
	channelHandler := make(chan Produce)

	go func() {
		vars := mux.Vars(r)
		code := vars["code"]
		var foundItem Produce
		for index, item := range Food {
			if item.ProduceCode == code {
				foundItem = Food[index]
				break
			}
		}
		channelHandler <- foundItem
	}()

	foundItem := <-channelHandler
	if foundItem.ProduceCode != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(foundItem)
}

func requestHandler() {
	foodRouter := mux.NewRouter().StrictSlash(true)
	foodRouter.HandleFunc("/produce", GetAllFoods)
	foodRouter.HandleFunc("/produce/{code}", GetFoodByCode)
	log.Fatal(http.ListenAndServe(":8000", foodRouter))
}

func main() {
	// Queue up the Food array
	Food = []Produce {
		{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.64"},
		{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}
	requestHandler()
}