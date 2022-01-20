package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Produce: structured list of produce items in market
type Produce struct {
	ProduceCode string
	Name        string
	UnitPrice   string
}

// Food: Array of Produce items
var Food []Produce

func main() {
	// Queue up the Food array
	Food = []Produce {
		Produce{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.64"},
		Produce{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		Produce{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		Produce{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}
}