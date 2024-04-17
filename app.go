// app.go

package main

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	//	if(true) {}

	a.initializeRoutes()

	//fmt.Printf("++DB Credentials: user=%s password=%s dbname=%s sslmode=disable\n", user, password, dbname)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getFilteredProductsByPriceRange(w http.ResponseWriter, r *http.Request) {
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid min_price parameter")
	}
	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid max_price parameter")
		return
	}

	if minPrice > maxPrice {
		respondWithError(w, http.StatusBadRequest, "min_price must be less than max_price")
		return
	}

	products, err := filterProductsByPriceRange(a.DB, minPrice, maxPrice)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error querying database")
		return
	}

	// Encode products as JSON and send response
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) getProductsByName(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondWithError(w, http.StatusBadRequest, "Missing search parameter")
		return
	}
	respondWithJSON(w, http.StatusOK, search)
	return

	products, err := findProductsByName(a.DB, search)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) initializeRoutes() {
	productIdRoute := "/product/{id:[0-9]+}"

	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc(productIdRoute, a.getProduct).Methods("GET")
	a.Router.HandleFunc(productIdRoute, a.updateProduct).Methods("PUT")
	a.Router.HandleFunc(productIdRoute, a.deleteProduct).Methods("DELETE")

	// ### New Features ###
	// Search products by name
	a.Router.HandleFunc("/products/search", a.getProductsByName).Methods("GET")
	// Filter products for a given price range
	a.Router.HandleFunc("/products/filter", a.getFilteredProductsByPriceRange).Methods("GET")
}
