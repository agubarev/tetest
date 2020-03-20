# Test homework for the TET company

A sample service project that imports the actual currency values from an [external feed](https://www.bank.lv/vk/ecb_rss.xml),
processes, stores and further serves via 2 public endpoints:

```
/api/v1/currency        -- returns a list of the latest known currency values
/api/v1/currency/:id    -- returns a historical list of currency values for a given currency ID (i.e.: USD)
```

## Getting Started

To get started, simply clone the repository and run `docker-compose up`

### Prerequisites

To run the project you must have docker and docker-compose installed on your system, please refer to the following links for more information.

* [docker installation guide](https://docs.docker.com/install/)
* [docker-compose installation guide](https://docs.docker.com/compose/install/)

### Installing

Clone the repository

```
git clone git@github.com:agubarev/tetest.git
cd tetest
```

run the tests

```
go test ./...
```

import the latest currency values (although it is done during the server startup)

```
docker exec -it app /bin/tetest import
```

see something like the following
```
➜  ~ docker exec -it app /bin/tetest import
2020-03-20T09:32:09.987Z	INFO	cmd/root.go:122	initializing database connection
2020-03-20T09:32:09.987Z	INFO	cmd/root.go:138	initializing default MySQL store
2020-03-20T09:32:09.987Z	INFO	cmd/root.go:147	initializing currency manager
2020-03-20T09:32:09.987Z	INFO	cmd/root.go:154	configuring main logger
2020-03-20T09:32:09.987Z	INFO	[currency]	cmd/import.go:32	importing currency data
2020-03-20T09:32:10.224Z	DEBUG	[currency]	currency/manager.go:135	parsing raw currency feed
```

run the whole thing via `docker-compose up`

```
docker-compose up -d
```

then, for example, to obtain the latest values of a specific currency, perform a GET request to `http://localhost:8080/api/v1/currency/USD` and get a sample response like the following (output copied from POSTMAN)
```
{
    "status_code": 200,
    "exec_time": 0.000312181,
    "payload": [
        {
            "id": "USD",
            "value": 1.0801,
            "pub_date": "2020-03-19T00:00:00Z",
            "created_at": "2020-03-20T07:21:31Z",
            "updated_at": "2020-03-20T07:25:37Z"
        },
        {
            "id": "USD",
            "value": 1.0934,
            "pub_date": "2020-03-18T00:00:00Z",
            "created_at": "2020-03-20T07:21:31Z",
            "updated_at": "2020-03-20T07:25:37Z"
        },
        {
            "id": "USD",
            "value": 1.0982,
            "pub_date": "2020-03-17T00:00:00Z",
            "created_at": "2020-03-20T07:21:31Z",
            "updated_at": "2020-03-20T07:25:37Z"
        },
        {
            "id": "USD",
            "value": 1.1157,
            "pub_date": "2020-03-16T00:00:00Z",
            "created_at": "2020-03-20T07:21:31Z",
            "updated_at": "2020-03-20T07:25:37Z"
        }
    ]
}
```

or GET `http://localhost:8080/api/v1/currency` to obtain the list of all the latest currency values

## Project Structure

Below is the file structure of this simple test project

```
➜  tetest git:(master) ✗ tree
.
├── bin                                             -- reserved for local testing and development
├── cmd                                             -- contains CLI commands
│   ├── import.go                                   -- represents the currency import command
│   ├── root.go
│   └── start.go                                    -- starts the server
├── db
│   └── baseline.sql                                -- the initial MySQL schema
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── internal 
│   ├── currency                                    -- business logic layer
│   │   ├── item.go
│   │   ├── manager.go
│   │   ├── store_cassandra.go
│   │   ├── store.go
│   │   ├── store_memory.go
│   │   ├── store_memory_test.go
│   │   ├── store_mysql.go
│   │   └── store_mysql_test.go
│   └── server                                      -- contains the server and endpoint implementation
│       ├── endpoints
│       │   ├── currency_get_by_id.go
│       │   ├── currency_get_latest.go
│       │   ├── currency_get_latest_test.go
│       │   └── endpoint.go
│       └── server.go
├── LICENSE
├── main.go
├── README.md
└── util                                            -- miscellaneous utilities which deserve their own space
    └── guard
        ├── guard.go
        └── guard_test.go

```
