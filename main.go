package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"

	"github.com/gorilla/mux"
	"github.com/asaskevich/govalidator"
)

// Produce: structured list of produce items in market
type Produce struct {
	ProduceCode string `json:"code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"price"`
}

// ItemValidator: Validation interface for required fields
type ItemValidator interface {
	Validate(r *http.Request) error
}

// Food: Array of Produce items
var Food []Produce

// We use Validate as a function to make sure user is entering required info
// ProduceCode should be 16 chars long alphanumeric and case insensitive with dashes 		 	  separating each group of 4 chars
// Name should be alphanumeric and case insensitive
// UnitPrice is a number with 2 decimal places
func (p Produce) Validate(r *http.Request) error {
	if len(p.ProduceCode) != 19 {
		return errors.New("invalid produce code")
	}
	chars := []rune(p.ProduceCode)

	if !govalidator.IsAlphanumeric(string(chars[0:3])) {
		return errors.New("invalid produce code")
	}
	return nil
}

// Interface method for item validation
func Validate(r *http.Request, v ItemValidator) error {
	return v.Validate(r)
}

// GetAllFoods: Displays all items in produce list
func GetAllFoods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /produce: GetAllFoods")
	json.NewEncoder(w).Encode(Food)
}

// GetFoodByCode: Displays particular item by product code
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

// AddFood: Method to create and add one or more new food items to database
func AddFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /produce: AddFood")

	reqBody, _ := ioutil.ReadAll(r.Body)
	x := bytes.TrimLeft(reqBody, " /t/r/n")
	isArray := len(x) > 0 && x[0] == '['

	itemsToAdd := make([]Produce, 0)
	if isArray {
		decoder := json.NewDecoder(bytes.NewBufferString(string(reqBody)))
		err := decoder.Decode(&itemsToAdd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, item := range itemsToAdd {
			err := Validate(r, item)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	} else {
		var addedItem Produce
		json.Unmarshal(reqBody, &addedItem)
		err := Validate(r, addedItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		itemsToAdd = append(itemsToAdd, addedItem)
	}

	channelHandler := make(chan []Produce)
	go func(items []Produce) {
			Food = append(Food, items...)
		channelHandler <- items
	}(itemsToAdd)
}

// DeleteFood: Method used to delete items from database
func DeleteFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /produce/{code}: DeleteFood")
	channelHandler := make(chan bool)

	go func() {
		vars := mux.Vars(r)
		code := vars["code"]
		var deletedItem bool = false
		for index, item := range Food {
			if item.ProduceCode == code {
				Food = append(Food[:index], Food[index + 1:]...)
				deletedItem = true
				break
			}
		}
		channelHandler <- deletedItem
	}()
}

// requestHandler: Will handle all router functions with mux
func requestHandler() {
	foodRouter := mux.NewRouter().StrictSlash(true)
	foodRouter.HandleFunc("/produce", GetAllFoods)
	foodRouter.HandleFunc("/produce/{code}", GetFoodByCode)
	foodRouter.HandleFunc("/produce", AddFood).Methods("POST")
	foodRouter.HandleFunc("/produce/{code}", DeleteFood).Methods("DELETE")
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