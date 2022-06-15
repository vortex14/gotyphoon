package cryptogen

import (
	"log"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

func GenerateSSHPair() (error, []byte, []byte) {
	bitSize := 4096
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil { return err, nil, nil}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil { return err, nil, nil}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	return err, publicKeyBytes, privateKeyBytes
}


// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil { return nil, err }

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil { return nil, err }

	log.Println("SSH Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns to the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil { return nil, err }

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("SSH Public key generated")
	return pubKeyBytes, nil
}



