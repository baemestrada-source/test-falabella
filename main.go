package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
//  "github.com/baemestrada-source/test-falabella/utils"
)

// BeerBox defines model for BeerBox.
type BeerBox struct {
	PriceTotal *float32 `json:"Price Total,omitempty"`
}

// BeerItem defines model for BeerItem.
type BeerItem struct {
	Brewery  string  `json:"Brewery"`
	Country  string  `json:"Country"`
	Currency string  `json:"Currency"`
	Id       int     `json:"Id"`
	Name     string  `json:"Name"`
	Price    float32 `json:"Price"`
}

var BeerArray []BeerItem

func GetBeersAll(w http.ResponseWriter, req *http.Request){
  json.NewEncoder(w).Encode(BeerArray)
}

func GetBeerID(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    vID, error := strconv.Atoi(params["beerID"])
    
    if error == nil {
        for _, item := range BeerArray {
            if item.Id == vID {
              json.NewEncoder(w).Encode(item)
              return
            }
          }
          json.NewEncoder(w).Encode(&BeerItem{})
    } else {
        log.Println(error)
    }

  }


func main() {
  router := mux.NewRouter()
  
  // adding example data
  BeerArray = append(BeerArray, BeerItem{Id: 1, Name:"Golden", Brewery:"kross",Country:"Chile", Price:10.5, Currency:"EUR" })
  BeerArray = append(BeerArray, BeerItem{Id: 2, Name:"Gallo",  Brewery:"Cerveceria CA",Country:"Guate", Price:5,    Currency:"QTZ" })

  // endpoints
  router.HandleFunc("/beers", GetBeersAll).Methods("GET")
  router.HandleFunc("/beers/{beerID}", GetBeerID).Methods("GET")

  log.Println("Servidor corriendo en el puerto 4000")

  log.Fatal(http.ListenAndServe(":4000", router))
}