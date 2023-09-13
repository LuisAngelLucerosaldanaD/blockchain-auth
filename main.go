package main

import (
	"blion-auth/api/grpc"
	"blion-auth/internal/env"
	"os"
)

func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)
	grpc.Start(c.App.Port)
}

/*func main() {
	clientPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
		return
	}

	serverPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
		return
	}
	clientPublicKey := clientPrivKey.PublicKey
	curve := elliptic.P256()
	x, _ := curve.ScalarMult(clientPublicKey.X, clientPublicKey.Y, serverPrivKey.D.Bytes())

	sharedSecret := x.Bytes()
	textC := encrypt([]byte("esto es un test"), sharedSecret)
	fmt.Println(string(textC))

	clientPublicKeyV2 := serverPrivKey.PublicKey
	curveV2 := elliptic.P256()
	xx, _ := curveV2.ScalarMult(clientPublicKeyV2.X, clientPublicKeyV2.Y, clientPrivKey.D.Bytes())
	sharedSecretV2 := xx.Bytes()
	fmt.Println(string(decrypt(textC, sharedSecretV2)))
}

func encrypt(message []byte, sharedKey []byte) []byte {
	block, _ := aes.NewCipher(sharedKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nonce, nonce, message, nil)
	return ciphertext
}

func decrypt(ciphertext []byte, sharedKey []byte) []byte {
	block, _ := aes.NewCipher(sharedKey)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}*/
