# Parrot API

# How to start DB?
- Run `make start-db` starts a Postgres docker and creates data tables.

# How to build?
- Run `make build` to compile the project directly on your OS, binary will be added to build directory, it requires go 1.16+.
- Run `make docker-build` to create a container, it requires docker.

# How to run an example?
- Run `make run` to start an instance of the API on port 8080.

# How to run tests?
- Run `make lint` to run linter.
- Run `make test` to execute unit tests.
- Run `make integration-test` to execute integration tests.
- Run `make load-test` to execute a load test (for this one you must have an instance of the API running).

# Postman
- Postman link: https://www.getpostman.com/collections/37f45b528a48e4cd38f3

# Endpoints

## CreateUser
### Description:
Method: POST
>```
>http://localhost:8080/user
>```
### Body (**raw**)

```json
{
    "email":"prueba@gmail.com",
    "fullName":"Carlos Flores",
    "password":"uno"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃


## CreateOrder
### Description:
Method: POST
>```
>http://localhost:8080/order
>```
### Headers

|Content-Type|Value|
|---|---|
|Authorization|Basic ajZyY3JjanJAZ21haWwuY29tOnVubw==|


### Body (**raw**)

```json
{
    "email":"prueba@gmail.com",
    "clientName":"Carl1os Flores",
    "price":1234,
    "products":[{
        "name":"uno",
        "price": 1,
        "description":"uno uno",
        "amount":5
    },{
        "name":"dos",
        "price": 2,
        "description":"uno uno",
        "amount":15
    },
    {
        "name":"tres",
        "price": 20,
        "description":"uno uno",
        "amount":15111
    }]
}
```

### 🔑 Authentication basic

|Param|value|Type|
|---|---|---|
|username|mail of the user|string|
|password|password selected by the user|string|


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃


## GenerateReport
### Description:
Method: POST
>```
>http://localhost:8080/report
>```
### Body (**raw**)

```json
{
    "from":"2015-01-28T17:41:52Z",
    "to": "2215-01-28T17:41:52Z"
}
```

### 🔑 Authentication basic

|Param|value|Type|
|---|---|---|
|username|mail of the user|string|
|password|password selected by the user|string|


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

_________________________________________________
Author: [Carlos Flores](https://github.com/hellerox)
