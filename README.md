# Getting Started

```shell
# download the project
git clone https://github.com/wondenge/go-genity.git

cd go-genity

# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

# run the RESTful API server
make run
```

At this time, you have a RESTful API server running at `http://127.0.0.1:5000`. It provides the following endpoints:

- `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
- `POST /v1/login`: authenticates a user and generates a JWT
- `GET /v1/genitys`: returns a paginated list of the genitys
- `GET /v1/genitys/:id`: returns the detailed information of an genity
- `POST /v1/genitys`: creates a new genity
- `PUT /v1/genitys/:id`: updates an existing genity
- `DELETE /v1/genitys/:id`: deletes an genity

Try the URL `http://localhost:5000/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

If you have `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following
more complex scenarios:

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:5000/v1/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the genity resources, such as: GET /v1/genitys
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:5000/v1/genitys
# should return a list of genity records in the JSON format
```

To use the project as a starting point of a real project whose package name is `github.com/abc/xyz`, do a global
replacement of the string `github.com/wondenge/go-genity` in all of project files with the string `github.com/abc/xyz`.

## Project Layout

The project uses the following project layout:

```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── genity            genity-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example,
the `genity` directory contains the application logic related with the genity feature.

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application
into a docker image. The docker container starts with the `cmd/server/entryscript.sh` script which reads
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `qa`, the application will be started with the `config/qa.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```
