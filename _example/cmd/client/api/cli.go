package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/googleapis/go-type-adapters/adapters"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	phonebookv1 "phonebook/pkg/phonebook/v1"
)

// version
var version = "client-v0.1.0"

// ClientCLI
type ClientCLI struct {
	app *cli.App
}

// NewClientCLI
func NewClientCLI() *ClientCLI {
	app := &cli.App{
		Name:     "pbclient",
		HelpName: "pbclient",
		Usage:    "phonebook client cli",
		Version:  version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "addr",
				Aliases:  []string{"a"},
				Usage:    "`address` to server",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			commandGetContact(),
			commandPutContact(),
			commandDeleteContact(),
		},
	}
	return &ClientCLI{
		app: app,
	}
}

// commandGetContact
func commandGetContact() *cli.Command {
	return &cli.Command{
		Name:  "find",
		Usage: "find contact",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "contact `name`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			ctxc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctxc, ctx.String("addr"), grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			store := phonebookv1.NewPhonebookStoreServiceClient(conn)

			resp, err := store.GetContact(context.Background(), &phonebookv1.GetContactRequest{
				FullName: ctx.String("name"),
			})
			if err != nil {
				return err
			}

			createdAt, err := adapters.ProtoDateTimeToTime(resp.Contact.CreatedAt)
			if err != nil {
				return err
			}

			fmt.Println("Name:", resp.Contact.FullName)
			fmt.Println("Email:", resp.Contact.Email)
			fmt.Println("Phone:", resp.Contact.Phone)
			fmt.Println("Create At:", createdAt.Format(time.ANSIC))

			return nil
		},
	}
}

// commandPutContact
func commandPutContact() *cli.Command {
	return &cli.Command{
		Name:  "save",
		Usage: "save a new contact",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "contact `name`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "email",
				Aliases:  []string{"e"},
				Usage:    "contact `email`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "phone",
				Aliases:  []string{"p"},
				Usage:    "contact `phone`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			ctxc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctxc, ctx.String("addr"), grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			store := phonebookv1.NewPhonebookStoreServiceClient(conn)

			createdAt, err := adapters.TimeToProtoDateTime(time.Now())
			if err != nil {
				return err
			}

			_, err = store.PutContact(context.Background(), &phonebookv1.PutContactRequest{
				Contact: &phonebookv1.Contact{
					FullName:  ctx.String("name"),
					Email:     ctx.String("email"),
					Phone:     ctx.String("phone"),
					CreatedAt: createdAt,
				},
			})
			if err != nil {
				return err
			}

			fmt.Println("Contact added")

			return nil
		},
	}
}

// commandDeleteContact
func commandDeleteContact() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete contact",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "contact `name`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			ctxc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctxc, ctx.String("addr"), grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			store := phonebookv1.NewPhonebookStoreServiceClient(conn)

			_, err = store.DeleteContact(context.Background(), &phonebookv1.DeleteContactRequest{
				FullName: ctx.String("name"),
			})
			if err != nil {
				return err
			}

			fmt.Println("Contact deleted")

			return nil
		},
	}
}

// // commandListContact
// func commandListContact() *cli.Command {
// 	return nil
// }

// Run
func (cli *ClientCLI) Run() error {
	return cli.app.Run(os.Args)
}
