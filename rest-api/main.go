package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Item an inventory item
type Item struct {
	ID    string  `json:"ID"`
	Name  string  `json:"Name"`
	Desc  string  `json:"Desc"`
	Price float64 `json:"Price"`
}

var inventory []Item

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Home Page</h1>")
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(inventory)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)

	inventory = append(inventory, item)

	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)

	for i, next := range inventory {
		if next.ID == id {
			inventory[i] = Item{
				ID:    id,
				Name:  item.Name,
				Desc:  item.Desc,
				Price: item.Price,
			}
			break
		}
	}

	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	for i, item := range inventory {
		if item.ID == id {
			inventory = append(inventory[:i], inventory[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(inventory)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/inventory", getInventory).Methods("GET")
	router.HandleFunc("/inventory", createItem).Methods("POST")
	router.HandleFunc("/inventory/{id:[0-9]+}", updateItem).Methods("PUT")
	router.HandleFunc("/inventory/{id:[0-9]+}", deleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	inventory = append(inventory, Item{
		ID:    "1",
		Name:  "Pumpkin",
		Desc:  "A pre-jack-o-lantern",
		Price: 5.99,
	})
	inventory = append(inventory, Item{
		ID:    "2",
		Name:  "Avalon Milk",
		Desc:  "Best milk in BC.",
		Price: 3.99,
	})
	handleRequests()
}
