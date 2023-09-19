package ciphers

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/Luzifer/go-openssl/v4"
)

func RsaPublicStringToRsaPublic(public string) *rsa.PublicKey {
	blockRsa, _ := pem.Decode([]byte(public))
	if blockRsa == nil {
		return nil
	}
	publicRsaPem, err := x509.ParsePKIXPublicKey(blockRsa.Bytes)
	if err != nil {
		return nil
	}

	publicRsa, ok := publicRsaPem.(*rsa.PublicKey)
	if !ok {
		return nil
	}
	return publicRsa
}

func RsaPrivateStringToRsaPrivate(public string) *rsa.PrivateKey {
	blockRsa, _ := pem.Decode([]byte(public))
	if blockRsa == nil {
		return nil
	}
	privateRsaPem, err := x509.ParsePKCS1PrivateKey(blockRsa.Bytes)
	if err != nil {
		return nil
	}

	return privateRsaPem
}

func GenerateKeyPairEcdsaX25519() (string, string, error) {
	newPrivateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	pemPrivateKey, err := EncodePrivateX25519(newPrivateKey)
	if err != nil {
		return "", "", err
	}

	pemPublicKey, err := EncodePublicX25519(newPrivateKey.PublicKey())
	if err != nil {
		return "", "", err
	}

	return pemPrivateKey, pemPublicKey, nil
}

func EncodePrivateX25519(privateKey *ecdh.PrivateKey) (string, error) {
	encoded, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})
	return string(pemEncoded), nil
}

func EncodePublicX25519(pubKey *ecdh.PublicKey) (string, error) {

	encoded, err := x509.MarshalPKIXPublicKey(pubKey)

	if err != nil {
		return "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	return string(pemEncodedPub), nil
}

func DecodePrivateX25519(pemEncodedPrivate string) (*ecdh.PrivateKey, error) {
	blockPrivate, _ := pem.Decode([]byte(pemEncodedPrivate))
	x509EncodedPrivate := blockPrivate.Bytes
	privateDecode, err := x509.ParsePKCS8PrivateKey(x509EncodedPrivate)
	if err != nil {
		return nil, err
	}
	private := privateDecode.(*ecdh.PrivateKey)
	return private, err
}

func DecodePublicX25519(pemEncodedPub string) (*ecdh.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdh.PublicKey)
	return publicKey, err
}

func GenerateKeyPairEcdsa() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	pemPrivateKey, err := EncodePrivate(privateKey)
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.PublicKey
	pemPublicKey, err := EncodePublic(&publicKey)
	if err != nil {
		return "", "", err
	}

	return pemPrivateKey, pemPublicKey, nil
}

func EncodePrivate(privateKey *ecdsa.PrivateKey) (string, error) {
	encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})
	return string(pemEncoded), nil
}

func EncodePublic(pubKey *ecdsa.PublicKey) (string, error) {

	encoded, err := x509.MarshalPKIXPublicKey(pubKey)

	if err != nil {
		return "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	return string(pemEncodedPub), nil
}

func DecodePrivate(pemEncodedPrivate string) (*ecdsa.PrivateKey, error) {
	blockPrivate, _ := pem.Decode([]byte(pemEncodedPrivate))
	x509EncodedPrivate := blockPrivate.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509EncodedPrivate)
	return privateKey, err
}

func DecodePublic(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey, err
}

func SignWithEcdsa(hash []byte, private ecdsa.PrivateKey) (string, error) {
	sign, err := ecdsa.SignASN1(rand.Reader, &private, hash)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func VerifySignWithEcdsa(hash []byte, pubKey ecdsa.PublicKey, sign []byte) (bool, error) {
	return ecdsa.VerifyASN1(&pubKey, hash, sign), nil
}

func StringToHashSha256(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func CipherDH(private ecdsa.PrivateKey, publicOtterKey ecdsa.PublicKey) (string, error) {
	curve := elliptic.P256()

	sharedKeyX, _ := curve.ScalarMult(publicOtterKey.X, publicOtterKey.Y, private.D.Bytes())
	message := []byte("Este es el contenido del mensaje a firmar")
	o := openssl.New()

	enc, err := o.EncryptBytes(sharedKeyX.String(), message, openssl.BytesToKeyMD5)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}
