# Yo-kai Watch 2 (3DS) replacement server
Includes both the authentication and secure servers

## Compiling

### Setup
Install [Go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads), then clone and enter the repository

```bash
$ git clone https://github.com/PretendoNetwork/yo-kai-watch-2
$ cd yo-kai-watch-2
```

### Compiling using `go`
To compile using Go, `go get` the required modules and then `go build` to your desired location. You may also want to tidy the go modules, though this is optional

```bash
$ go get -u
$ go mod tidy
$ go build -o build/yo-kai-watch-2
```

The server is now built to `build/yo-kai-watch-2`

When compiling with only Go, the authentication servers build string is not automatically set. This should not cause any issues with gameplay, but it means that the server build will not be visible in any packet dumps or logs a title may produce

To compile the servers with the authentication server build string, add `-ldflags "-X 'main.serverBuildString=BUILD_STRING_HERE'"` to the build command, or use `make` to compile the server

### Compiling using `make`
Compiling using `make` will read the local `.git` directory to create a dynamic authentication server build string, based on your repositories remote origin and current commit. It will also use the current folders name as the executables name

Install `make` onto your system (this varies by OS), and run `make` while inside the repository

```bash
$ make
```

The server is now built to `build/yo-kai-watch-2` with the authentication server build string already set

## Configuration
All configuration options are handled via environment variables

`.env` files are supported

| Name                                 | Description                                                                                                         | Required                                      |
|--------------------------------------|---------------------------------------------------------------------------------------------------------------------|-----------------------------------------------|
| `PN_YKW2_POSTGRES_URI`               | Fully qualified URI to your Postgres server (Example `postgres://username:password@localhost/YKW2?sslmode=disable`) | Yes                                           |
| `PN_YKW2_AUTHENTICATION_SERVER_PORT` | Port for the authentication server                                                                                  | Yes                                           |
| `PN_YKW2_SECURE_SERVER_HOST`         | Host name for the secure server (should point to the same address as the authentication server)                     | Yes                                           |
| `PN_YKW2_SECURE_SERVER_PORT`         | Port for the secure server                                                                                          | Yes                                           |
| `PN_YKW2_ACCOUNT_GRPC_HOST`          | Host name for your account server gRPC service                                                                      | Yes                                           |
| `PN_YKW2_ACCOUNT_GRPC_PORT`          | Port for your account server gRPC service                                                                           | Yes                                           |
| `PN_YKW2_ACCOUNT_GRPC_API_KEY`       | API key for your account server gRPC service                                                                        | No (Assumed to be an open gRPC API)           |
| `PN_YKW2_FRIENDS_GRPC_HOST`          | Host name for your friends server gRPC service                                                                      | Yes                                           |
| `PN_YKW2_FRIENDS_GRPC_PORT`          | Port for your friends server gRPC service                                                                           | Yes                                           |
| `PN_YKW2_FRIENDS_GRPC_API_KEY`       | API key for your friends server gRPC service                                                                        | No (Assumed to be an open gRPC API)           |
| `PN_YKW2_AES_KEY`                    | AES key from the account server                                                                                     | Yes                                           |
