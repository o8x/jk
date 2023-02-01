package cert

import (
	"crypto/x509"
	"fmt"
	"os"
)

type Folder struct {
	Crt     []byte `json:"crt"`
	Key     []byte `json:"key"`
	Csr     []byte `json:"csr"`
	CrtFile string `json:"crt_file"`
	CsrFile string `json:"csr_file"`
	KeyFile string `json:"key_file"`
}

func AutoFolder(domain string) (*Folder, error) {
	if _, err := MakeCA(domain); err != nil {
		return nil, err
	}

	if _, err := MakeCertFromCSR(domain); err != nil {
		return nil, err
	}

	folder := &Folder{
		CrtFile: fmt.Sprintf("/var/tmp/%s.crt", domain),
		CsrFile: fmt.Sprintf("/var/tmp/%s.csr", domain),
		KeyFile: fmt.Sprintf("/var/tmp/%s.key", domain),
	}

	if err := folder.ParseFile(); err != nil {
		return nil, err
	}

	return folder, nil
}

func (f *Folder) GetCertPool() (*x509.CertPool, error) {
	if err := f.ParseFile(); err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(f.Crt)
	return pool, nil
}

func (f *Folder) ParseFile() error {
	var err error
	if f.KeyFile != "" {
		f.Key, err = os.ReadFile(f.KeyFile)
		if err != nil {
			return err
		}
	}

	if f.CrtFile != "" {
		f.Crt, err = os.ReadFile(f.CrtFile)
		if err != nil {
			return err
		}
	}

	if f.CsrFile != "" {
		f.Csr, err = os.ReadFile(f.CsrFile)
		if err != nil {
			return err
		}
	}

	return nil
}
