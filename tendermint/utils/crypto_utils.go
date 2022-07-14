package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
)

type KeyPair struct {
	PublicKey [] byte
	PrivateKey [] byte
}

func Hash(msg string) []byte {
	h := sha256.New()
	h.Write([]byte("hello world\n"))
	out := h.Sum(nil)
	return out
}

func GenerateNewKeyPair() (*KeyPair, error) {
	pubKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		ErrorStdOutLogger.Printf("Generate Key Failed: %s", err)
		return nil, err
	}
	keyPair := KeyPair{PublicKey: pubKey, PrivateKey: privateKey}
	return &keyPair, err
}

func SignData(privateKey []byte, data []byte) []byte {
	return ed25519.Sign(privateKey, data)
}

func VerifySignature(pubKey, data, signature []byte) bool {
	return ed25519.Verify(pubKey, data, signature)
}
