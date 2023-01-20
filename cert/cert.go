package cert

import (
	"fmt"
	"os/exec"
	"strings"
)

// openssl req -text -in localhost.csr
// openssl x509 -noout -text -in localhost.crt
// openssl req -new -subj "/C=US/ST=Utah/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr
// openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt

func MakeCA(subject string) error {
	args := strings.Fields(fmt.Sprintf(`req -new -subj /C=US/ST=Utah/CN=%s -newkey rsa:2048 -nodes -keyout %s.key -out %s.csr`, subject, subject, subject))
	return exec.Command("openssl", args...).Run()
}

func MakeCertFromCSR(domain string) error {
	args := strings.Fields(fmt.Sprintf(`x509 -req -days 365 -in %s.csr -signkey %s.key -out %s.crt`, domain, domain, domain))
	return exec.Command("openssl", args...).Run()
}
