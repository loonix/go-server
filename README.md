# go-server

A very simple api written in go.
It uses fake data variable to store books.
Uses gin web framework to deal with http requests.

## Run go app
```shell
go run main.go
```

### How it looks in the terminal
<img width="554" alt="image" src="https://user-images.githubusercontent.com/3384277/233748792-783ed49e-0826-46fe-ab52-42eb43409202.png">


### API Endpoints

**Add** a new book from a fake data file 

```shell
curl -X POST -d @addbook.json -H "Content-Type: application/json" localhost:8080/books
```
**Get** all books

<img width="1040" alt="image" src="https://user-images.githubusercontent.com/3384277/233748724-382050f5-363d-46a0-85b1-644b14db1a8c.png">

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
