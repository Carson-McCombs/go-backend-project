#About
* Developed with Go 1.23, using Gorilla Mux as a HTTP Router, and Docker as a packaging / delivery tool. See the summary.txt file for more.


# User Guide:
* Note, permission request might pop up for things like accessing ports and building files, please accept this requests for the program to function as intended.


# Prerequisites:

**(You might not need all of these depending on execution method)**

* Install Go: https://go.dev/doc/install

* Install Docker: https://docs.docker.com/engine/install/

* Install Gorilla Mux: ( after Go has been installed and within a Cmd / Terminal ) 

```go get -u github.com/gorilla/mux```

## Without Docker:

* Install Go and Gorilla Mux

* Open CMD / Terminal

* Navigate to desired directory

```
git clone https://github.com/Carson-McCombs/go-backend-project/
cd go-backend-project
go run main.go
```


## For Build to Execution with Docker:

* Ensure that Go, Docker, and Gorilla Mux is installed

* Open Cmd / Terminal

* Navigate to desired directory

```
git clone https://github.com/Carson-McCombs/go-backend-project/
cd go-backend-project
docker build -t go-fetch-backend . 
docker run -p 8000:8000 go-fetch-backend
```

## For Docker Image ( without Go or Gorilla Mux )

* Ensure that Docker is installed and setup

* Open Cmd

* Navigate to desired folder

```
docker pull carsonmccombs/go-fetch-backend
docker run -p 8000:8000 go-fetch-backend
```
