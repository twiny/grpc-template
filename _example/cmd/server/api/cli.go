package api

import (
	"log"
	"net"
	"os"

	phonebookv1 "phonebook/pkg/phonebook/v1"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// version
var version = "server-v0.1.0"

// CLI
type ServerCLI struct {
	app *cli.App
}

// NewServerCLI
func NewServerCLI() *ServerCLI {
	app := &cli.App{
		Name:     "phonebook",
		HelpName: "phonebook",
		Usage:    "phonebook cli",
		Version:  version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "store",
				Aliases:  []string{"s"},
				Usage:    "`path` to store.db",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "server `port`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			return start(ctx.String("store"), ctx.String("port"))
		},
	}
	return &ServerCLI{
		app: app,
	}
}

// Run
func (cli *ServerCLI) Run() error {
	return cli.app.Run(os.Args)
}

// start
func start(path, port string) error {
	store, err := NewStore(path)
	if err != nil {
		return err
	}
	defer store.Close()

	server := grpc.NewServer()

	phonebookv1.RegisterPhonebookStoreServiceServer(server, store)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Println("server running at: ", port)

	return server.Serve(listener)
}
