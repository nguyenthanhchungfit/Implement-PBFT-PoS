package utils

import "testing"

func TestSignature(t *testing.T) {
	rawMsg := "Hello, world."
	rawData := Hash(rawMsg)
	pubKey, privKey, err := GenerateNewKeyPair()
	if err == nil{
		sign := SignData(privKey, rawData)
		verified := VerifySignature(pubKey, rawData, sign)
		if !verified {
			t.Error("Signature is failed")
		}
	}
}
