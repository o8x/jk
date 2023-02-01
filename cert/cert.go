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

func MakeCA(subject string) ([]byte, error) {
	args := strings.Fields(fmt.Sprintf(`req -new -subj /SAN=ALEX/C=US/ST=Utah/CN=%s -newkey rsa:2048 -nodes -keyout %s.key -out %s.csr`, subject, subject, subject))

	path, err := exec.LookPath("openssl")
	if err != nil {
		return nil, err
	}

	cmd := exec.Cmd{
		Path: path,
		Args: args,
		Dir:  "/var/tmp",
	}
	return cmd.CombinedOutput()
}

func MakeCertFromCSR(domain string) ([]byte, error) {
	args := strings.Fields(fmt.Sprintf(`x509 -req -days 365 -in %s.csr -signkey %s.key -out %s.crt`, domain, domain, domain))

	path, err := exec.LookPath("openssl")
	if err != nil {
		return nil, err
	}

	cmd := exec.Cmd{
		Path: path,
		Args: args,
		Dir:  "/var/tmp",
	}
	return cmd.CombinedOutput()
}
