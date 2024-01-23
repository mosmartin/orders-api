package handlers

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...creating order 🚧")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...List all orders 🚧")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...Getting order by ID 🚧")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...Updating order by ID 🚧")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...Deleting order by ID 🚧")
}
