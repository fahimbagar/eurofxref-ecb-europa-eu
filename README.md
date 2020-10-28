# Euro Forex - European Central Bank
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-89%25-brightgreen.svg?longCache=true&style=flat)</a>

## About the Project
This project is to show historical rates from https://www.ecb.europa.eu using [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and use only the standard library other than database connection.

## Prerequisites
- [Golang](https://golang.org/)
```shell script
go version go1.15.2 linux/amd64
```
- [Go SQLite3](https://github.com/mattn/go-sqlite3) 
- Optional: [Docker](https://www.docker.com)

## Usage
- Build Golang files and run the executable
```shell script
make
./app/eurofxref-ecb
```

- Or, run using docker
```shell script
make docker
docker run --name eurofxref-ecb -p 8282:8282 --rm fahimbagar/eurofxref-ecb:1.0
```

## Test
Run test coverage
```shell script
make test
```

## List of Endpoints API
- Get forex latest rates: http://localhost:8282/rates/latest

- Get forex rates at certain dates (format: yyyy-mm-dd): http://localhost:8282/rates/2020-10-27

- Analyze forex rates: http://localhost:8282/rates/analyze
