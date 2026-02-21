package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	return crypto.Sign(hash[:], privateKey)
}

func VerifySignature(publicKey *ecdsa.PublicKey, message, signature []byte) bool {
	hash := sha256.Sum256(message)
	return crypto.VerifySignature(curveToBytes(publicKey), hash[:], signature)
}

func VerifySignatureRaw(publicKey []byte, message, signature []byte) bool {
	hash := sha256.Sum256(message)
	return crypto.VerifySignature(publicKey, hash[:], signature)
}

func curveToBytes(curve *elliptic.Curve) []byte {
	return curve.Params().Name[:]
}

func MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalECPrivateKey(key)
}

func UnmarshalPrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(data)
}

func PublicKeyToBytes(key *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(key.Curve, key.X, key.Y)
}

func PublicKeyFromBytes(data []byte, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, data)
	if x == nil {
		return nil, fmt.Errorf("failed to unmarshal public key")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

func PrivateKeyToBytes(key *ecdsa.PrivateKey) []byte {
	return x509.MarshalECPrivateKey(key)
}

func PrivateKeyFromBytes(data []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(data)
}

func EncodePEM(key []byte, keyType string) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  keyType,
		Bytes: key,
	})
}

func DecodePEM(data []byte) ([]byte, string, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, "", fmt.Errorf("no PEM block found")
	}
	return block.Bytes, block.Type, nil
}

func ComputeHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func VerifyHash(data, hash []byte) bool {
	computed := sha256.Sum256(data)
	return string(computed[:]) == string(hash)
}
