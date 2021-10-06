## phonebook grpc server/client

## Generating packages
run `make generate`

## Example Commands

### start server
`go run _example/cmd/server/main.go --store "tmp/store.db" --port 8813`

### add contact
`go run _example/cmd/client/main.go --addr "127.0.0.1:8813" save --name "james" --email "james@home.com" --phone "333-111-2123"`

### find contact
`go run _example/cmd/client/main.go --addr "127.0.0.1:8813" find --name "mike"`

### delete contact
`go run _example/cmd/client/main.go --addr "127.0.0.1:8813" delete --name "mike"`

### list all contact
`go run _example/cmd/client/main.go --addr "127.0.0.1:8813" -h`