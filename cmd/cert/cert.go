package main

import (
	"github.com/o8x/jk/v2/args"
	"github.com/o8x/jk/v2/cert"
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:  "cert maker",
			Usage: "go run github.com/o8x/jk/v2/cmd/cert -domain localhost",
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
		if _, err := cert.MakeCA(d); err != nil {
			a.PrintErrorExit(err)
		}

		if _, err := cert.MakeCertFromCSR(d); err != nil {
			a.PrintErrorExit(err)
		}
	}
}
