package main

import (
	"github.com/o8x/jk/args"
	"github.com/o8x/jk/cert"
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:  "cert maker",
			Usage: "go run github.com/o8x/jk/cmd/cert -domain localhost",
		},
		Flags: []*args.Flag{
			{
				Name:        []string{"-domain", "-d"},
				Description: "cert subject for domain",
				Required:    true,
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintHelpExit(err)
	}

	for _, d := range a.Gets("domain") {
		if err := cert.MakeCA(d); err != nil {
			a.PrintErrorExit(err)
		}

		if err := cert.MakeCertFromCSR(d); err != nil {
			a.PrintErrorExit(err)
		}
	}
}
