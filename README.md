##


## Commands

### add contact
`go run cmd/client/main.go --addr "127.0.0.1:8813" save --name "mike" --email "mike@home.com" --phone "888-333-5553"`

### find contact
`go run cmd/client/main.go --addr "127.0.0.1:8813" find --name "mike"`

### delete contact
`go run cmd/client/main.go --addr "127.0.0.1:8813" delete --name "mike"`