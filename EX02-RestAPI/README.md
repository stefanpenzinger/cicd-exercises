# Exercise 02
Accompanying Chapter  Chapters „Programming in GO“ and „Microservice Architectures“

## Description
In this exercise you will implement a small rest-api (microservice) in Go
- Follow the [instructions](https://semaphoreci.com/community/tutorials/building-andtesting-a-rest-api-in-go-with-gorilla-mux-and-postgresql)
  - work until the headline „Setting Up Continuous Integration with Semaphore“
- Prerequisites:
  - postgreSQL installed (user/role „postgres“)
  - Github Account
- The code and additional hints can be found [here](https://github.com/mrckurz/go-mux)

## Goal
Make sure to **UNDERSTAND** what is happening in the code.<br>
Add additional functionality to the code.<br>
Implemented functionality:
- Added route to filter products by a price range
- Added route to search products by their name

## Start Microservice
Run the following in order to start the microservice
```shell
./start_go_rest_api.sh
```

Run the following in order to test the microservice
```shell
./test_go_rest_api.sh
```