Development Choices:

Programming Language: Go / Golang

Reasoning: I have prior experience with multiple languages / web-frameworks that could be used here ( such as Go, PHP, Django / Python ), but I know Go is built for web services, scales well, and is optimized for concurrency. I also know that Fetch uses Go for some of its backend servers and services.

Development Environment: Visual Studio Code

Reasoning: It is open source, has plenty of support for many languages and frameworks, and I am used to it.

Other Tools: Docker

Reasoning: I've used Docker before and it works well for ensuring a quick deployment, testing, and compatibility. It also means that the user does not have to worry about setting anything else up or troubleshooting dependencies since it is run in pre-configured internal VM.

Libraries / Dependencies: Gorilla Mux

Reasoning: I have worked with Gorilla Mux before and know it is used within the industry as an HTTP router for Go web servers.


User Guide:
*Note, permission request might pop up for things like accessing ports and building files, please accept this requests for the program to function as intended.

Without Docker:
*Install Go
```
go get -u github.com/gorilla/mux
```
*Open CMD / Terminal
*Navigate to desired directory
```
git clone https://github.com/Carson-McCombs/go-backend-project/
cd go-backend-project
go run main.go
```


For Build to Execution with Docker:

*Ensure that Go, Docker, and Gorilla Mux is installed
*Open Cmd / Terminal
*Navigate to desired directory
```
git clone https://github.com/Carson-McCombs/go-backend-project/
cd go-backend-project
docker build -t go-fetch-backend . 
docker run -p 8000:8000 go-fetch-backend
```

For Docker Image to Execution:

*Ensure that Docker is installed and setup
*Open Cmd
*Navigate to desired folder
```
docker pull carsonmccombs/go-fetch-backend
docker run -p 8000:8000 go-fetch-backend
```
