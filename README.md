# go-server

A very simple api written in go.
It uses fake data variable to store books.
Uses gin web framework to deal with http requests.

## Run go app
```shell
go run main.go
```

### API Endpoints

**Add** a new book from a fake data file 

```shell
curl -X POST -d @addbook.json -H "Content-Type: application/json" localhost:8080/books
```
**Get** all books
```shell
curl localhost:8080/books
```

**Get** a book **by id**
```shell
curl localhost:8080/books/1
```

**Update** a book by id
```shell
curl -X PUT -d @updatebook.json -H "Content-Type: application/json" localhost:8080/books/1
```

**Delete** a book by id
```shell
curl -X DELETE localhost:8080/books/1
```