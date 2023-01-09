# quiz

## working with application

* start server `go run main.go serve`
* start client `go run main.go play`

## working with buf

* linting proto files `buf lint`
* formatting proto files `buf format -w`
* check breaking changes `buf breaking --against ".git#branch=main"`
