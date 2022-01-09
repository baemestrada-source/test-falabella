package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
  "io/ioutil"
)

//ResponseTC estructura para leer y decodificar el json de la API de tipo de cambio
type ResponseTC struct {
  Success   bool   `json:"success"` //si devuelve lo correcto es true
  Quotes map[string]float32 `json:"quotes"` //devuelve un array en forma de mapa de mis claves de monedas y valores
}

// BeerBox defines model for BeerBox.
type BeerBox struct {
	PriceTotal float32 `json:"Price Total,omitempty"`
}

// BeerItem defines model for BeerItem.
type BeerItem struct {
	Id       int     `json:"Id"`
	Name     string  `json:"Name"`
  Brewery  string  `json:"Brewery"`
	Country  string  `json:"Country"`
	Price    float32 `json:"Price"`
  Currency string  `json:"Currency"`
}

//BeerArray estructura donde se llena la informacion en memoria
var BeerArray []BeerItem

//Acces Key que me devolvio la API para tipo de cambio
const access_key = "a88c3af250d29d5d2c77e71066c27a92" //clave para la API de tipo de cambio

//TpcIn aqui se guardara el tipo de cambio USD versus la moneda que se desea 
var TpcIn float32 

//TpcDb aqui se guardara el tipo de cambio USD versus la moneda que tiene el precio en mi data
var TpcDb float32 

//searchBeers Lista todas las cervezas
func searchBeers(w http.ResponseWriter, req *http.Request){
  //Simplemente devuelvo todo lo que tenga en el array
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(BeerArray)

}

//searchBeerById Busca una cerveza por su Id
func searchBeerById(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)

    // la variable que obtengo de el path
    beerID, err := strconv.Atoi(params["beerID"])

    // si no hay errores inicio
    if err == nil {
        for _, item := range BeerArray {
            if item.Id == beerID {
              //el ID fue encontrado lo retorno en JSON
              w.Header().Set("Content-Type", "application/json")
              w.WriteHeader(http.StatusOK)
              json.NewEncoder(w).Encode(item)
              return
            }
          }
    } 
    http.Error(w,"El Id de la cerveza no existe", http.StatusNotFound)
  }

  //addBeers Ingresa una nueva cerveza
  func addBeers(w http.ResponseWriter, req *http.Request){
    var beer BeerItem
    err := json.NewDecoder(req.Body).Decode(&beer)

    //Unicamente entra al if si no hay ningun error con el body
    if err == nil {
      for _, item := range BeerArray {
        if item.Id == beer.Id {
          http.Error(w,"El ID de la cerveza ya existe", 409)
          return
        }
      }
      //aqui si no encontro ningun error regreso el item y respondo que todo esta bien status 201
      BeerArray = append(BeerArray, beer)

      //retorno el mensaje al servidor para indicar que la cerveza fue creada exitosamente
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusCreated)
      json.NewEncoder(w).Encode("Cerveza creada")

      } else { 
      //error 400 ya que el body tiene error 
      http.Error(w,"Request invalida", http.StatusBadRequest)
    }
  }

  //boxBeerPriceById Obtiene el precio de una caja de cerveza por su Id
  func boxBeerPriceById(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)

    query := req.URL.Query()
    currency := query.Get("currency")

    //Cantidad de cervezas a comprar (valor por defecto 6), la cantidad viene en string la paso a float32
    quantity, err := strconv.ParseFloat(query.Get("quantity"),32)
    if err != nil {
      quantity = 6
    }

    //paso a entero para comparar 
    beerID, err := strconv.Atoi(params["beerID"])
    if err == nil {
        for _, item := range BeerArray {
            if item.Id == beerID {
              //consulto la API con la variable que ya tenia definida
              response, err := http.Get("http://api.currencylayer.com/live?access_key="+access_key+"&format=1")      
              if err != nil {
                  http.Error(w, "Ocurrio un error al leer la API de tipo de cambio", http.StatusBadRequest)
                  return
              }
              
              responseData, err := ioutil.ReadAll(response.Body)
              if err != nil {
                http.Error(w, "Ocurrio un error al leer la API de tipo de cambio", http.StatusBadRequest)
                return
              }

              //log.Println(string(responseData))
              
              var responseObject ResponseTC
              json.Unmarshal(responseData, &responseObject)

              if responseObject.Success {
                for k, v := range responseObject.Quotes {
                  // Busca el tipo de cambio con respecto al USD de la moneda con la que se pagara
                    if k[3:6] == currency { 
                      TpcIn = v
                    }
                    // Busca el tipo de cambio con respecto al USD a la moneda que esta en la data
                    if k[3:6] == item.Currency {  
                      TpcDb = v
                    }
                }
                //valor homologado utilizando los valores de la API
                Valor_H :=  ( item.Price / TpcDb ) * TpcIn
                //log.Println(TpcIn,TpcDb, Valor_H)


                //Calculo el precio total ya con la moneda que corresponde
                ValPTot:=Valor_H  * float32(quantity)  

                // Lleno mi estructura de respuesta
                BeerBox:=BeerBox{PriceTotal:ValPTot} 

                //respondo la API con la estructura
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(&BeerBox)
              }
          
              return
            }
          }
    } 
    http.Error(w,"El Id de la cerveza no existe", http.StatusNotFound)

  }

func main() {
  router := mux.NewRouter()
  
  // Informacion de ejemplo de entrada
  BeerArray = append(BeerArray, BeerItem{Id: 1, Name:"Golden", Brewery:"kross",Country:"Chile", Price:10.5, Currency:"USD" })
  BeerArray = append(BeerArray, BeerItem{Id: 2, Name:"Gallo",  Brewery:"Cerveceria CA",Country:"Guate", Price:5,    Currency:"GTQ" })

  // endpoints
  router.HandleFunc("/beers", searchBeers).Methods("GET")
  router.HandleFunc("/beers", addBeers).Methods("POST")
  router.HandleFunc("/beers/{beerID}", searchBeerById).Methods("GET")
  router.HandleFunc("/beers/{beerID}/boxprice", boxBeerPriceById).Methods("GET")

  //log 
  log.Println("Servidor corriendo en el puerto 4000")

  //escucha del servidor
  log.Fatal(http.ListenAndServe(":4000", router))
}