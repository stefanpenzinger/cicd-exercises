# Exercise 02 + 03 ![example workflow](https://github.com/stefanpenzinger/cicd-exercises/actions/workflows/go.yml/badge.svg)
Accompanying Chapter  Chapters „Programming in GO“ and „Microservice Architectures“

## Exercise 02
### Description
In this exercise you will implement a small rest-api (microservice) in Go
- Follow the [instructions](https://semaphoreci.com/community/tutorials/building-andtesting-a-rest-api-in-go-with-gorilla-mux-and-postgresql)
  - work until the headline „Setting Up Continuous Integration with Semaphore“
- Prerequisites:
  - postgreSQL installed (user/role „postgres“)
  - Github Account
- The code and additional hints can be found [here](https://github.com/mrckurz/go-mux)

### Goal
Make sure to **UNDERSTAND** what is happening in the code.<br>
Add additional functionality to the code.<br>
Implemented functionality:
- Added route to filter products by a price range
- Added route to search products by their name

### Start Microservice
Run the following in order to start the microservice
```shell
./start_go_rest_api.sh
```

Run the following in order to test the microservice
```shell
./test_go_rest_api.sh
```

## Exercise 03
### Description
Implement the following CI steps:
- Create a file go.yml in your repository (folder .github/workflows/) and add the required elements.
- Trigger a build by a code change and commit to the repository.
- Add build status badge to the README.md
- Integrate SonarCloud (https://sonarcloud.io) in your process - everytime a new commit is done into the repo, a new Sonar-run should be triggered