# GetGround Party App

The GetGround Part App allows managing the part guestlist.

Since the capacity and the number of tables are variable (not specified in requirement), we allow the table's creation on demand.

When a guest wants to add his name to the list, it must provide the "number" of the table (as he sees it on a layout picture).

If the table does not exist, we create a table with the chosen number, default table size, and available spots.

If the table does exist, we check if the remaining available seats are enough for his party.


The table size is defined during the app initialization as an environment variable (default size 12).

## How to run:

```
make docker-up # Initilizalizes db and app containers
```

## Use Cases

- Add a guest to the list:

```
Request:

POST localhost:3000/guest_list/john
{
    "table": 1,
    "accompanying_guests": 3
}

Response: 

201 Created
{
  "Name": "john"
}

```

- Get guestlist:
```
Request:

GET localhost:3000/guest_list

Response:

200 OK
{
  "guests": [
    {
      "name": "alesr",
      "table": 1,
      "accompanying_guests": 3
    }
  ]
}

```

- Welcome guest (guest arrives)
```
Request:

PUT localhost:3000/guests/john
{
    "accompanying_guests": 3
}

Response:

200 OK
{
  "name": "john"
}
```

- Get arrived guests:
```
Request:

GET localhost:3000/guests

Response:

200 OK
{
  "guests": [
    {
      "name": "john",
      "accompanying_guests": 3,
      "time_arrived": "2021-11-23T12:23:31Z"
    }
  ]
}
```

- Goodbye guest (guest leaves):
```
Request:

DELETE localhost:3000/guests/john

Response:

200 OK
```

- Get empty seats:
```
Request:

GET localhost:3000/seats_empty

Response:

200 OK
{
  "empty_seats": 20
}
```

## Code structure

```
.
├── cmd
│   └── party
│       └── main.go  // Initialization
├── docker-compose.yaml
├── Dockerfile
├── internal
│   ├── app // REST application 
│   │   ├── app.go
│   │   └── partyctrl
│   │       ├── partyctrl.go  // Party HTTP Handlers (transport layer)
│   │       └── partyctrl_test.go
│   └── pkg
│       └── party
│           ├── errors.go
│           ├── mock.go
│           ├── models.go
│           ├── model_test.go
│           ├── party.go // Party Service (Domain layer)
│           ├── party_test.go
│           └── repository
│               ├── mock.go
│               ├── mysql.go
│               ├── mysql_test.go
│               └── repository.go // Part repository (storage layer)
├── Makefile
├── pkg
│   └── database
│       └── database.go
└── README.md
```

[Standard Go Project Layout](https://github.com/golang-standards/project-layout)
