package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	apple_mdm "github.com/fleetdm/fleet/v4/server/mdm/apple"
)

// from MicroMDM
const (
	rsaPrivateKeyPEMBlockType         = "RSA PRIVATE KEY"
	pushCertificatePrivateKeyFilename = "PushCertificatePrivateKey.key"
	mdmcertdir                        = "mdm-certificates"
)

func encryptedKey(key *rsa.PrivateKey, password []byte) ([]byte, error) {
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPEMBlock, err := x509.EncryptPEMBlock(rand.Reader, rsaPrivateKeyPEMBlockType, privBytes, password, x509.PEMCipher3DES)
	if err != nil {
		return nil, err
	}

	out := pem.EncodeToMemory(privPEMBlock)
	return out, nil
}

func run() error {
	flEmail := flag.String("email", "", "Email address to use in CSR Subject.")
	flOrg := flag.String("org", "", "Organization to use in the CSR Subject.")
	flPKeyPass := flag.String("password", "", "Password to encrypt/read the RSA key.")
	flKeyPath := flag.String("private-key", filepath.Join(mdmcertdir, pushCertificatePrivateKeyFilename), "Path to the push certificate private key. A new RSA key will be created at this path.")
	flag.Parse()

	if *flEmail == "" {
		return errors.New("-email must be set")
	}
	if *flOrg == "" {
		return errors.New("-org must be set")
	}
	if *flPKeyPass == "" {
		return errors.New("-password must be set")
	}

	csr, key, err := apple_mdm.GenerateAPNSCSRKey(*flEmail, *flOrg)
	if err != nil {
		return fmt.Errorf("could not generate CSR: %w", err)
	}

	pemKey, err := encryptedKey(key, []byte(*flPKeyPass))
	if err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	if err := os.WriteFile(*flKeyPath, pemKey, 0600); err != nil {
		return fmt.Errorf("could not write key: %w", err)
	}
	fmt.Println("wrote private key to", *flKeyPath)

	if err := apple_mdm.GetSignedAPNSCSR(http.DefaultClient, csr); err != nil {
		return fmt.Errorf("could not submit CSR: %w", err)
	}

	fmt.Println("CSR submitted to Fleet's server. Check inbox at", *flEmail, "for email from Fleet.")
	return nil
}

func main() {
	if err := run(); err != nil {
		flag.Usage()
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
