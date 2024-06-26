// main_test.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		"postgres", //os.Getenv("APP_DB_USERNAME"),
		"password", //os.Getenv("APP_DB_PASSWORD"),
		"postgres") //os.Getenv("APP_DB_NAME"))

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Print("\t" + string(body) + "\n")

}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Print("\t" + string(body) + "\n")

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Print("\t" + string(body) + "\n")
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Print("\t" + string(body) + "\n")

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestFilterProducts(t *testing.T) {
	created := 10
	expected := 4

	minPrice := 10
	maxPrice := minPrice * expected

	clearTable()
	addProducts(created)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/products/filter?min_price=%d&max_price=%d", minPrice, maxPrice), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var products []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &products)

	if len(products) != expected {
		t.Errorf("Expected %d products. Got %d", expected, len(products))
	}
}

func TestFilterProductsMissingQueryParams(t *testing.T) {
	missingQueryParams, _ := http.NewRequest("GET", "/products/filter", nil)
	response := executeRequest(missingQueryParams)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	missingMinPriceParam, _ := http.NewRequest("GET", "/products/filter?max_price=40", nil)
	response = executeRequest(missingMinPriceParam)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	missingMaxPriceParam, _ := http.NewRequest("GET", "/products/filter?min_price=10", nil)
	response = executeRequest(missingMaxPriceParam)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestFilterProductsMinPriceBiggerThanMaxPrice(t *testing.T) {
	minPriceBiggerThanMaxPrice, _ := http.NewRequest("GET", "/products/filter?min_price=40&max_price=10", nil)
	response := executeRequest(minPriceBiggerThanMaxPrice)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestSearchProducts(t *testing.T) {
	searchTerm := "Product"
	expected := 11

	clearTable()
	addProducts(expected)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/products/search?term=%s", searchTerm), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var products []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &products)

	if len(products) != expected {
		t.Errorf("Expected %d products. Got %d", expected, len(products))
	}

	expected = 2
	searchTerm = "1"

	req, _ = http.NewRequest("GET", fmt.Sprintf("/products/search?term=%s", searchTerm), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &products)

	if len(products) != expected {
		t.Errorf("Expected %d products. Got %d", expected, len(products))
	}
}

func TestSearchProductsMissingQueryParam(t *testing.T) {
	missingQueryParam, _ := http.NewRequest("GET", "/products/search", nil)
	response := executeRequest(missingQueryParam)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
