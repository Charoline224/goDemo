package net

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type dollars float32

type database map[string]dollars

func main() {
	db := database{"shoes": 2, "cloth": 1}
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/query", http.HandlerFunc(db.query))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s", item, price)
	}
	return
}
func (db database) query(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "price: %s", price)
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	_, ok := db[item]
	if ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "the %s already exisit", item)
		return
	}
	price := req.URL.Query().Get("price")
	priceF, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "price is invalid: %q \n", price)
		return
	}
	db[item] = dollars(priceF)
	fmt.Fprintf(w, "create item: %q\n", item)
	fmt.Fprintf(os.Stdout, "item create: %q\n", db)
}
