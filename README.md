[![CI Linter pipeline](https://github.com/harness/gitness/actions/workflows/ci-lint.yml/badge.svg)](https://github.com/harness/gitness/actions/workflows/ci-lint.yml)
[![CodeQL](https://github.com/harness/gitness/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/harness/gitness/actions/workflows/codeql-analysis.yml)
[![](https://img.shields.io/badge/go-%3E%3D%201.19-green)](#)
# Pre-Requisites

Install the latest stable version of Node and Go version 1.19 or higher, and then install the below Go programs. Ensure the GOPATH [bin directory](https://go.dev/doc/gopath_code#GOPATH) is added to your PATH.

```bash
$ make all
```

Setup github access token required for UI dependencies:
```bash
$ yarn setup-github-registry
```

# Build

Build the user interface:

```bash
$ pushd web
$ yarn install
$ yarn run build
$ popd
```

Build the server and command line tools:

```bash
# STANDALONE
$ make build
```

# Test

Execute the unit tests:

```bash
$ make test
```

# Run

This project supports all operating systems and architectures supported by Go.  This means you can build and run the system on your machine; docker containers are not required for local development and testing.

Start the server at `localhost:3000`

```bash
# STANDALONE
./gitness server .local.env
```

# User Interface

This project includes a simple user interface for interacting with the system. When you run the application, you can access the user interface by navigating to `http://localhost:3000` in your browser.

# Swagger

This project includes a swagger specification. When you run the application, you can access the swagger specification by navigating to `http://localhost:3000/swagger` in your browser.

# CLI
This project includes simple command line tools for interacting with the system. Please remember that you must start the server before you can execute commands.

Register a new user:

```bash
$ ./gitness register
```

> NOTE: A user `admin` (pw: `changeit`) gets created by default.


Login to the application:

```bash
$ ./gitness login
```

Logout from the application:

```bash
$ ./gitness logout
```

View your account details:

```bash
$ ./gitness user self
```

Generate a personal access token:

```bash
$ ./gitness user pat $NAME $LIFETIME_IN_S
```

Debug and output http responses from the server:

```bash
$ DEBUG=true ./gitness user self
```

View all commands:

```bash
$ ./gitness --help
```

# REST API
Please refer to the swagger for the specification of our rest API.

For testing, it's simplest to execute operations as the default user `admin` using a PAT:
```bash
# LOGIN (user: admin, pw: changeit)
$ ./gitness login

# GENERATE PAT (1 YEAR VALIDITY)
$ ./gitness user pat mypat 2592000
```

The command outputs a valid PAT that has been granted full access as the user.
The token can then be send as part of the `Authorization` header with Postman or curl:

```bash
$ curl http://localhost:3000/api/v1/user \
-H "Authorization: Bearer $TOKEN"
```